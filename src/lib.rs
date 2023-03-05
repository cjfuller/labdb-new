use hyper::{client::HttpConnector, Client};
use hyper_rustls::HttpsConnector;

pub mod auth;
pub mod db;
pub mod env;
pub mod errors;
pub mod models;
pub mod proxy;
pub mod search;

pub(crate) fn http_client() -> Client<HttpsConnector<HttpConnector>> {
    let https = hyper_rustls::HttpsConnectorBuilder::new()
        .with_webpki_roots()
        .https_or_http()
        .enable_http1()
        .build();

    hyper::Client::builder().build(https)
}
