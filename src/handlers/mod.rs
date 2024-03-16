use axum::response::{Html, IntoResponse};

// GET requests
pub async fn handler_startpage() -> impl IntoResponse {
    Html("Hello, startpage!")
}

pub async fn handler_viewer() -> impl IntoResponse {
    Html("Hello, viewer!")
}

pub async fn handler_image_request() -> impl IntoResponse {
    Html("Hello, image_request!")
}

pub async fn handler_fullscreen() -> impl IntoResponse {
    Html("Hello, fullscreen!")
}

// POST requests
pub async fn handler_user_image_data() -> impl IntoResponse {
    Html("Hello, user_image_data!")
}
