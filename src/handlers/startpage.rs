use askama::Template;
use axum::response::{Html, IntoResponse};

#[derive(Template)]
#[template(path = "startpage.html")]

struct StartpageTemplate {
    number_of_images: u32,
}

pub async fn generate() -> impl IntoResponse {
    let startpage = StartpageTemplate {
        number_of_images: 800,
    };

    Html(startpage.render().unwrap())
}
