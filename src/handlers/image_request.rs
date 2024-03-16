use askama::Template;
use axum::response::{Html, IntoResponse};

#[derive(Template)]
#[template(path = "image_request.html")]

struct ImageRequestTemplate {
    image_id: u32,  // This is unique key from database
    image_num: u16, // same example this would be: 32
    time: String,
    trial_num: u8, // same example would be: 3
}

pub async fn generate() -> impl IntoResponse {
    let image_request_page = ImageRequestTemplate {
        image_id: 123, // This is unique key from database
        image_num: 32, // same example this would be: 32
        time: "2024-05-13 13:51:32".to_string(),
        trial_num: 3,
    };

    Html(image_request_page.render().unwrap())
}
