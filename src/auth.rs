use std::collections::HashMap;

use axum::async_trait;
use axum::extract::{FromRequestParts, Query, State};
use axum::http::request::Parts;
use axum::http::{HeaderMap, Request, StatusCode};
use axum::middleware::Next;
use axum::response::{IntoResponse, Redirect, Response};
use axum_sessions::extractors::{ReadableSession, WritableSession};
use hmac::Mac;
use hyper::Method;
use serde::Deserialize;
use sha2::Sha256;
use sqlx::{PgConnection, PgPool};

use crate::errors::{Error, Result};
use crate::http_client;

const APP_ID: &str = "146923434465-alq7iagpanjvoag20smuirj0ivdtfldk.apps.googleusercontent.com";

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
    let Some(token) = params.get("token") else {
        return Ok((StatusCode::BAD_REQUEST, "No auth token provided").into_response());
    };
    match get_verified_identity(token).await? {
        None => Ok((StatusCode::FORBIDDEN, "Forbidden").into_response()),
        Some(verified) => {
            session.insert("user_id", verified)?;
            Ok(Redirect::to("/").into_response())
        }
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
