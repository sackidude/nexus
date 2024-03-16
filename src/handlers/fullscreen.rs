use askama::Template;
use axum::response::{Html, IntoResponse};

#[derive(Template)]
#[template(path = "fullscreen.html")]

struct FullscreenTemplate {
    trial_num: u16,
    start_timestamp: String,
    yeast_amount: f32,
    sugar_amount: f32,
    stirring: bool,
}

pub async fn generate() -> impl IntoResponse {
    let fullscreen = FullscreenTemplate {
        trial_num: 3,
        start_timestamp: "2023-20-23 23:11".to_string(),
        yeast_amount: 4.0,
        sugar_amount: 20.0,
        stirring: true,
    };

    Html(fullscreen.render().unwrap())
}
