use axum::{
    routing::{get, get_service, post},
    Router,
};
use tower_http::services::ServeDir;

mod error;
mod handlers;

#[tokio::main]
async fn main() {
    println!("Started");
    let routes_all = Router::new()
        .merge(routes_htmx())
        .fallback_service(routes_static());

    let listener = tokio::net::TcpListener::bind("127.0.0.1:8080")
        .await
        .unwrap();
    axum::serve(listener, routes_all.into_make_service())
        .await
        .unwrap();
}

fn routes_htmx() -> Router {
    Router::new()
        // GET requests
        .route("/", get(handlers::handler_startpage))
        .route("/startpage", get(handlers::handler_startpage))
        .route("/viewer", get(handlers::handler_viewer))
        .route("/image-request", get(handlers::handler_image_request))
        .route("/fullscreen", get(handlers::handler_fullscreen))
        // POST requests
        .route("/user-image-data", post(handlers::handler_user_image_data))
}

fn routes_static() -> Router {
    Router::new().nest_service("/", get_service(ServeDir::new("./static")))
}
