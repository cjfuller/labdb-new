use axum::middleware::Next;
use axum::response::Response;
use axum::routing::{any, get, post};
use axum::Router;
use axum_sessions::{async_session, SessionLayer};
use hyper::Request;
use labdb::auth::{post_api_verify, require_auth};
use labdb::db::create_pool;
use labdb::env::Env;
use labdb::proxy::{proxy_public, proxy_with_auth};
use labdb::search;

fn setup_logger() -> Result<(), fern::InitError> {
    fern::Dispatch::new()
        .format(|out, message, record| {
            out.finish(format_args!(
                "[{}][{}][{}] {}",
                chrono::Local::now().format("%Y-%m-%dT%H:%M:%S"),
                record.target(),
                record.level(),
                message
            ))
        })
        .level(log::LevelFilter::Info)
        .chain(std::io::stdout())
        .apply()?;
    Ok(())
}

async fn request_logging<B>(request: Request<B>, next: Next<B>) -> Response {
    let start = std::time::Instant::now();
    let method = request.method().to_string();
    let path = request
        .uri()
        .path_and_query()
        .map(|it| it.to_string())
        .unwrap_or_else(|| "<path missing>".to_string());
    let resp = next.run(request).await;
    let dur = start.elapsed();
    let status = resp.status().as_u16();

    log::info!("{method} {path} :: {status} ({}ms)", dur.as_millis());
    resp
}

#[tokio::main]
async fn main() {
    setup_logger().unwrap();

    let store = async_session::CookieStore::new();
    let secret = labdb::env::SECRET_TOKEN.as_bytes();
    let session_layer = SessionLayer::new(store, secret)
        .with_cookie_name("labdb_session")
        .with_session_ttl(Some(std::time::Duration::from_secs(60 * 60 * 24 * 14)))
        .with_secure(true);

    let pool = create_pool().await.unwrap();

    // build our application with a single route
    let app = Router::new()
        // Auth-protected routes.
        .route("/search", get(search::search_handler))
        .layer(axum::middleware::from_fn_with_state(
            pool.clone(),
            require_auth,
        ))
        // Routes below this point are open access!
        .route("/", any(proxy_public))
        .route("/_s/*rest", any(proxy_public))
        .route("/api/verify", post(post_api_verify))
        // End open access routes.
        .fallback(proxy_with_auth)
        .layer(session_layer)
        .layer(axum::middleware::from_fn(request_logging))
        .with_state(pool);

    let listen_address = if Env::dev() {
        "127.0.0.1:3000".parse().unwrap()
    } else {
        let port = std::env::var("PORT").unwrap_or_else(|_| "3000".to_string());
        format!("0.0.0.0:{port}").parse().unwrap()
    };
    axum::Server::bind(&listen_address)
        .serve(app.into_make_service())
        .await
        .unwrap();
}
