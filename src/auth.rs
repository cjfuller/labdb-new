use std::collections::{HashMap, HashSet};
use std::time::Instant;

use axum::async_trait;
use axum::extract::{FromRequestParts, Query, State};
use axum::http::request::Parts;
use axum::http::{HeaderMap, Request, StatusCode};
use axum::middleware::Next;
use axum::response::{IntoResponse, Redirect, Response};
use axum_sessions::extractors::{ReadableSession, WritableSession};
use base64::Engine;
use hmac::Mac;
use hyper::Method;
use jwt_simple::prelude::{RS256PublicKey, RSAPublicKeyLike, VerificationOptions};
use jwt_simple::token::Token;
use once_cell::sync::Lazy;
use regex::Regex;
use serde::{Deserialize, Serialize};
use sha2::Sha256;
use sqlx::{PgConnection, PgPool};
use tokio::sync::RwLock;

use crate::errors::{Error, Result};
use crate::http_client;

const APP_ID: &str = "146923434465-alq7iagpanjvoag20smuirj0ivdtfldk.apps.googleusercontent.com";

type ParsedKeys = HashMap<String, RS256PublicKey>;

static G_KEYS: Lazy<RwLock<(Instant, ParsedKeys)>> =
    Lazy::new(|| RwLock::new((Instant::now(), HashMap::new())));
const G_KEYS_URL: &str = "https://www.googleapis.com/oauth2/v3/certs";

#[derive(Clone, Debug, Deserialize)]
struct Key {
    kid: String,
    alg: String,
    n: String,
    e: String,
}

#[derive(Clone, Debug, Deserialize)]
struct KeysResponse {
    keys: Vec<Key>,
}

async fn fetch_and_cache_google_keys() -> Result<ParsedKeys> {
    let resp = reqwest::get(G_KEYS_URL).await?;
    let cache_re = Regex::new(r"max-age=(\d+)").unwrap();
    let max_age: u64 = resp
        .headers()
        .get("Cache-Control")
        .and_then(|it| it.to_str().ok())
        .and_then(|it| cache_re.captures(it))
        .and_then(|caps| caps.get(1).map(|it| it.as_str()))
        .unwrap_or("0")
        .parse()
        .unwrap_or(0);
    let expiry: Instant = Instant::now() + std::time::Duration::from_secs(max_age);
    let keys: KeysResponse = resp.json().await?;
    let engine = base64::engine::general_purpose::URL_SAFE_NO_PAD;
    let mut parsed_keys: ParsedKeys = HashMap::new();
    for key in keys.keys {
        if &key.alg != "RS256" {
            log::error!("Ignoring google key with algorithm {}", key.alg);
            continue;
        }
        let Ok(n) = engine.decode(&key.n) else { continue };

        let Ok(e) = engine.decode(&key.e) else { continue };
        match RS256PublicKey::from_components(&n, &e) {
            Ok(parsed) => {
                let parsed = parsed.with_key_id(&key.kid);
                parsed_keys.insert(key.kid, parsed);
            }
            Err(e) => log::error!("Unable to parse RS256 public key: {e}"),
        }
    }
    *G_KEYS.write().await = (expiry, parsed_keys.clone());
    Ok(parsed_keys)
}

async fn google_key(key_id: &str) -> Result<RS256PublicKey> {
    let (exp, keys) = { G_KEYS.read().await.clone() };
    let keys = if exp < Instant::now() {
        fetch_and_cache_google_keys().await?
    } else {
        keys
    };
    keys.get(key_id)
        .ok_or_else(|| Error::str(format!("No key with ID {key_id}")))
        .cloned()
}

#[derive(Debug, Clone, Serialize, Deserialize)]
struct Claims {
    email: String,
    email_verified: bool,
}

async fn get_verified_token_email(token_string: &str) -> Result<String> {
    let meta = Token::decode_metadata(token_string)
        .map_err(|e| Error::str(format!("Failed to parse token: {e}")))?;
    let kid = meta
        .key_id()
        .ok_or_else(|| Error::str("Token did not include a key id."))?;
    let signing_key = google_key(kid).await?;
    let mut opts = VerificationOptions::default();
    let mut audiences = HashSet::new();
    audiences.insert(APP_ID.into());
    opts.allowed_audiences = Some(audiences);
    let claims: Claims = signing_key
        .verify_token(token_string, Some(opts))
        .map_err(|e| Error::str(format!("Error verifying token: {e}")))?
        .custom;
    if claims.email_verified {
        Ok(claims.email)
    } else {
        Err(Error::str("Unverified e-mail"))
    }
}

pub fn add_auth_headers(user_id: &str, headers: &mut HeaderMap) {
    let time = chrono::Utc::now().format("%FT%T").to_string();
    let mut mac = hmac::Hmac::<Sha256>::new_from_slice(crate::env::SIGNING_KEY.as_bytes()).unwrap();
    mac.update(format!("{user_id}{time}").as_bytes());
    let result = hex::encode(mac.finalize().into_bytes());
    headers.insert("X-LabDB-UserId", user_id.parse().unwrap());
    headers.insert("X-LabDB-Signature", result.parse().unwrap());
    headers.insert("X-LabDB-Signature-Timestamp", time.parse().unwrap());
}

#[derive(Debug)]
pub struct User {
    user_id: String,
}

#[derive(Clone, Copy, Debug)]
pub enum Permission {
    Read,
    Write,
    Admin,
}

impl Permission {
    pub fn for_method(method: &Method) -> Self {
        match *method {
            Method::GET | Method::HEAD | Method::OPTIONS => Permission::Read,
            _ => Permission::Write,
        }
    }
}

impl User {
    pub async fn verify_in_db(&self, conn: &mut PgConnection, permission: Permission) -> bool {
        let Some(user) = crate::models::User::find_by_email(conn, &self.user_id).await else {
            return false};
        match permission {
            Permission::Read => user.auth_read.unwrap_or_default(),
            Permission::Write => user.auth_write.unwrap_or_default(),
            Permission::Admin => user.auth_admin.unwrap_or_default(),
        }
    }
}

pub trait UserID {
    fn id(&self) -> Option<&str>;
}

impl UserID for User {
    fn id(&self) -> Option<&str> {
        Some(&self.user_id)
    }
}

impl UserID for Option<User> {
    fn id(&self) -> Option<&str> {
        match self {
            Some(u) => u.id(),
            None => None,
        }
    }
}

#[async_trait]
impl<S: Send + Sync> FromRequestParts<S> for User {
    type Rejection = StatusCode;

    async fn from_request_parts(
        parts: &mut Parts,
        state: &S,
    ) -> std::result::Result<Self, Self::Rejection> {
        let Ok(session) = ReadableSession::from_request_parts(parts, state).await else {
            return Err(StatusCode::FORBIDDEN);
        };
        session
            .get("user_id")
            .map(|user_id| User { user_id })
            .ok_or(StatusCode::FORBIDDEN)
    }
}

#[derive(Clone, Debug, Deserialize)]
struct IdentityResponse {
    aud: String,
    email_verified: String,
    email: String,
}

async fn get_verified_identity(token: &str) -> Result<Option<String>> {
    let verify_request = Request::builder()
        .method("POST")
        .uri(format!(
            "https://www.googleapis.com/oauth2/v3/tokeninfo?id_token={}",
            urlencoding::encode(token)
        ))
        .header("Content-Length", 0)
        .body(hyper::Body::empty())
        .unwrap();

    let mut verify_response = http_client().request(verify_request).await?;
    let id_response: IdentityResponse = if verify_response.status().is_success() {
        serde_json::from_slice(&hyper::body::to_bytes(verify_response.body_mut()).await?)?
    } else {
        log::error!(
            "Got error from identity verifier: {} {}",
            verify_response.status().as_u16(),
            String::from_utf8(
                hyper::body::to_bytes(verify_response.body_mut())
                    .await?
                    .to_vec()
            )
            .unwrap_or_else(|_| "<Non-text response body>".into())
        );
        return Err(Error::str("Got error from identity verifier"));
    };
    if id_response.aud.contains(APP_ID) && &id_response.email_verified == "true" {
        Ok(Some(id_response.email))
    } else {
        Ok(None)
    }
}

// POST /api/verify
pub async fn post_api_verify(
    Query(params): Query<HashMap<String, String>>,
    mut session: WritableSession,
) -> Result<Response> {
    if let Some(token) = params.get("jwt") {
        match get_verified_token_email(token).await {
            Err(e) => Ok((StatusCode::FORBIDDEN, format!("Forbidden: {e}")).into_response()),
            Ok(verified) => {
                session.insert("user_id", verified)?;
                Ok(Redirect::to("/").into_response())
            }
        }
    } else if let Some(token) = params.get("token") {
        match get_verified_identity(token).await? {
            None => Ok((StatusCode::FORBIDDEN, "Forbidden").into_response()),
            Some(verified) => {
                session.insert("user_id", verified)?;
                Ok(Redirect::to("/").into_response())
            }
        }
    } else {
        return Ok((StatusCode::BAD_REQUEST, "No auth token provided").into_response());
    }
}

pub async fn require_auth<B>(
    State(pool): State<PgPool>,
    request: Request<B>,
    next: Next<B>,
) -> Response {
    let (mut parts, body) = request.into_parts();
    let mut conn = pool.acquire().await.unwrap();
    match User::from_request_parts(&mut parts, &()).await {
        Ok(u)
            if u.verify_in_db(&mut conn, Permission::for_method(&parts.method))
                .await =>
        {
            let next_req = Request::from_parts(parts, body);
            next.run(next_req).await
        }
        _ => (StatusCode::FORBIDDEN, "Forbidden").into_response(),
    }
}
