use std::collections::HashMap;

use axum::extract::{Host, Query, State};
use axum::http::uri::PathAndQuery;
use axum::response::Response;
use hyper::{Method, Request, StatusCode, Uri};
use regex::Regex;
use serde::Serialize;
use sqlx::{PgConnection, PgPool};

use crate::auth::User;
use crate::errors::Result;
use crate::model_eval;
use crate::models::{owner_search_query, search_query, QueryResult};
use crate::proxy::proxy_with_auth;

fn hyper_body_respond(status: StatusCode, msg: &str) -> Result<Response<hyper::Body>> {
    Ok(Response::builder()
        .status(status)
        .body(msg.to_string().into())?)
}

#[derive(Debug, Serialize)]
struct BackendSearchBody {
    items: Vec<QueryResult>,
    term: String,
}

pub async fn search_handler<B>(
    Query(params): Query<HashMap<String, String>>,
    State(pool): State<PgPool>,
    user: User,
    host: Host,
    request: Request<B>,
) -> Result<Response<hyper::Body>> {
    let Some(term) = params.get("term") else {
        return hyper_body_respond(StatusCode::BAD_REQUEST, "Invalid search query");
    };
    let include_sequence = params.get("seq") == Some(&("1".to_string()));
    let person: Option<&str> = params.get("person").map(|it| it.as_ref());
    let Ok(types): std::result::Result<Vec<String>, _> = serde_json::from_str(params.get("types").map(|it| it.as_ref()).unwrap_or_default()) else {
        return hyper_body_respond(StatusCode::BAD_REQUEST, "Invalid search query");
    };

    let mut conn = pool.acquire().await?;
    let results = search(&mut conn, term, include_sequence, person, &types).await?;
    let request_body = BackendSearchBody {
        items: results,
        term: term.clone(),
    };
    let serialized = serde_json::to_string(&request_body).unwrap();
    let body: hyper::Body = serialized.into();

    let (mut parts, _) = request.into_parts();
    let mut uri_parts = parts.uri.clone().into_parts();
    uri_parts.path_and_query = Some(PathAndQuery::from_static("/search_result"));
    parts.uri = Uri::from_parts(uri_parts)?;
    parts.method = Method::POST;
    parts
        .headers
        .insert("Content-Type", "application/json".parse().unwrap());
    let backend_request = Request::from_parts(parts, body);
    let resp = proxy_with_auth(host, user, backend_request).await?;
    Ok(resp)
}

pub async fn search(
    conn: &mut PgConnection,
    term: &str,
    include_sequence: bool,
    person: Option<&str>,
    types: &[String],
) -> Result<Vec<QueryResult>> {
    let case_insensitive = term.starts_with('/') && term.ends_with("/i");
    let is_regexp = term.starts_with('/') && (term.ends_with("/i") || term.ends_with('/'));

    let mut norm_term: String = if is_regexp {
        let temp = Regex::new("^/").unwrap().replace(term, "");
        Regex::new("/i?$").unwrap().replace(&temp, "").to_string()
    } else {
        format!("^{}$", term.replace('*', ".*"))
    };

    if case_insensitive {
        norm_term = format!("(?i){norm_term}");
    }

    let search_re = Regex::new(&norm_term)?;

    let mut results = vec![];
    for model_type in types {
        let mts: &str = model_type.as_ref();
        model_eval!(mts, M, {
            if let Some(p) = person {
                results.append(&mut owner_search_query::<M>(conn, p).await?);
            } else {
                results.append(&mut search_query::<M>(conn, &search_re, include_sequence).await?);
            }
            Ok(())
        })?
    }

    Ok(results)
}
