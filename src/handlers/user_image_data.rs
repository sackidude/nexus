use axum::response::IntoResponse;

pub async fn generate() -> impl IntoResponse {
    super::image_request::generate().await
}
