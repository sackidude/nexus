use sqlx::Row;
use std::sync::{Arc, Mutex};

use sqlx::postgres::PgPool;

use crate::{
    error::{Error, Result},
    helpers::get_database_credentials,
};

#[derive(Clone)]
pub struct ModelController {
    pub pool: PgPool,
}

impl ModelController {
    pub async fn new() -> Result<Self> {
        sqlx::migrate!("db/migrations");

        let credentials = get_database_credentials()?;
        let url = format!(
            "postgres://{}:{}@localhost:5432/nexus",
            credentials.username, credentials.password
        );
        let pool = PgPool::connect(&url)
            .await
            .map_err(|_| Error::DatabaseConnectionFail)?;

        Ok(ModelController { pool: pool })
    }
}

impl ModelController {
    pub async fn get_image_count(&self) -> u32 {
        let res = sqlx::query("SELECT COUNT(*) AS img_count FROM images")
            .fetch_one(&self.pool)
            .await
            .expect("Failed to query database");

        let count: i32 = res.get("img_count");
        count.try_into().unwrap()
    }
}
