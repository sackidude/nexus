use askama::Template;
use axum::response;

#[derive(Template)]
#[template(path = "startpage.html")]

struct StartpageTemplate {
    number_of_images: u32,
}

use response as r;
pub async fn generate() -> impl r::IntoResponse {
    let startpage = StartpageTemplate {
        number_of_images: 800,
    };

    r::Html(startpage.render().unwrap())
}
