use askama::Template;
use axum::response::{Html, IntoResponse};

#[derive(Template)]
#[template(path = "viewer.html")]

struct ViewerTemplate {
    trial_nums: Vec<u16>,
    chart_data: String,
}

pub async fn generate() -> impl IntoResponse {
    let viewerpage = ViewerTemplate {
        trial_nums: vec![3, 4, 5],
        chart_data: "alert(1)".to_string(),
    };

    Html(viewerpage.render().unwrap())
}
