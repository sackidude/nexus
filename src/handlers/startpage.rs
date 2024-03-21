use sqlx::Row;
use std::sync::Arc;

use askama::Template;
use axum::{
    extract::State,
    response::{Html, IntoResponse},
};

use crate::AppState;

#[derive(Template)]
#[template(path = "startpage.html")]

struct StartpageTemplate {
    number_of_images: u32,
}
#[axum::debug_handler]
pub async fn generate(State(app_state): State<Arc<AppState>>) -> impl IntoResponse {
    let img_count: i64 = sqlx::query("SELECT COUNT(*) AS img_count FROM images")
        .fetch_one(&app_state.db)
        .await
        .unwrap()
        .get("img_count");

    let startpage = StartpageTemplate {
        number_of_images: img_count.try_into().unwrap(), // Infalliable because count can't be negative.
    };

    Html(startpage.render().unwrap())
}
