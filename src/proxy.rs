use std::env;

use axum::extract::Host;
use axum::http::uri::Scheme;
use axum::http::{Request, Uri};
use axum::response::Response;
use hyper::StatusCode;

use crate::auth::{User, UserID};
use crate::env::Env;
use crate::errors::{Error, Result};
use crate::http_client;

const DEV_PROXY_TARGET: &str = "localhost:3001";
const PROXY_SUFFIX: &str = "-backend.labdb.io";

fn backend_uri(uri: Uri, host: &str) -> Result<Uri> {
    let host = if Env::dev() {
        env::var("PROXY_TARGET").unwrap_or_else(|_| DEV_PROXY_TARGET.into())
    } else {
        if !host.ends_with(".labdb.io") {
            return Err(Error::str(
                "Can only proxy requests on .labdb.io subdomains.",
            ));
        }
        host.replace(".labdb.io", PROXY_SUFFIX)
    };
    let mut parts = uri.into_parts();
    parts.authority = Some(host.parse()?);
    parts.scheme = Some(parts.scheme.unwrap_or_else(|| {
        if Env::dev() {
            Scheme::HTTP
        } else {
            Scheme::HTTPS
        }
    }));
    Ok(Uri::try_from(parts)?)
}

fn set_up_backend_request<T>(host: &str, req: &mut Request<T>, maybe_current_user: Option<&str>) {
    let headers_to_remove = req
        .headers()
        .iter()
        .filter(|(k, _)| k.as_str().to_ascii_lowercase().starts_with("cf-"))
        .map(|(k, _)| k.clone())
        .collect::<Vec<_>>();

    for header in headers_to_remove {
        req.headers_mut().remove(header);
    }
    req.headers_mut().remove("X-Forwarded-For");
    req.headers_mut().remove("Host");
    req.headers_mut().remove("X-Forwarded-Host");
    req.headers_mut().remove("Forwarded");
    req.headers_mut()
        .insert("X-Labdb-Forwarded", "true".parse().unwrap());
    if let Some(user) = maybe_current_user {
        crate::auth::add_auth_headers(user, req.headers_mut());
    }
    let new_uri = backend_uri(req.uri().clone(), host).unwrap();
    log::info!("Proxying {:?} to {:?}", req.uri(), new_uri);
    if let Some(new_host) = new_uri.host() {
        req.headers_mut().insert("Host", new_host.parse().unwrap());
    }
    *req.uri_mut() = new_uri;
}

async fn proxy_internal(
    host: &str,
    user: Option<User>,
    mut request: Request<hyper::Body>,
) -> Result<Response<hyper::Body>> {
    if let Some(h) = request.headers().get("X-Labdb-Forwarded") {
        if matches!(h.to_str(), Ok("true") | Err(..)) {
            return Ok(Response::builder()
                .status(StatusCode::BAD_REQUEST)
                .body(hyper::Body::from("Stuck in a recursive proxy loop."))?);
        }
    }

    set_up_backend_request(host, &mut request, user.id());
    Ok(http_client().request(request).await?)
}

pub async fn proxy_with_auth(
    Host(host): Host,
    user: User,
    request: Request<hyper::Body>,
) -> Result<Response<hyper::Body>> {
    proxy_internal(&host, Some(user), request).await
}

pub async fn proxy_public(
    Host(host): Host,
    request: Request<hyper::Body>,
) -> Result<Response<hyper::Body>> {
    proxy_internal(&host, None, request).await
}
