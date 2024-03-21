use std::sync::Arc;

use axum::{
    debug_handler,
    routing::{get, get_service, post},
    Router,
};
use error::Result;
use sqlx::{database, PgPool};
use tower_http::services::ServeDir;

use crate::{error::Error, helpers::get_database_credentials};

mod error;
mod handlers;
mod helpers;
mod model;

pub struct AppState {
    db: PgPool,
}

#[tokio::main]
async fn main() -> Result<()> {
    dotenv::dotenv().ok();
    let database_url = std::env::var("DATABASE_URL").expect("DATABASE_URL must be set");

    let pool = PgPool::connect(&database_url)
        .await
        .map_err(|_| Error::DatabaseConnectionFail)?;
    let app =
        routes_htmx(Arc::new(AppState { db: pool.clone() })).fallback_service(routes_static());

    let listener = tokio::net::TcpListener::bind("127.0.0.1:8080")
        .await
        .unwrap();
    axum::serve(listener, app).await.unwrap();

    Ok(())
}

fn routes_htmx(app_state: Arc<AppState>) -> Router {
    use handlers as h;
    Router::new()
        // GET requests
        .route("/", get(h::startpage::generate))
        .route("/startpage", get(h::startpage::generate))
        // .route("/viewer", get(h::viewer::generate))
        // .route("/image-request", get(h::image_request::generate))
        // .route("/fullscreen", get(h::fullscreen::generate))
        // // POST requests
        // .route("/user-image-data", post(h::user_image_data::generate))
        .with_state(app_state)
}

fn routes_static() -> Router {
    Router::new().nest_service("/", get_service(ServeDir::new("./static")))
}
