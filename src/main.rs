use axum::{
    routing::{get, get_service, post},
    Router,
};
use tower_http::services::ServeDir;

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
}

fn routes_static() -> Router {
    Router::new().nest_service("/", get_service(ServeDir::new("./static")))
}
