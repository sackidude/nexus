use crate::error::{Error, Result};

pub struct DatabaseCredentials {
    pub username: String,
    pub password: String,
}

pub fn get_database_credentials() -> Result<DatabaseCredentials> {
    let username = dotenv::var("SQL_USERNAME").map_err(|_| Error::DatabaseCredentialsNotFound)?;
    let password = dotenv::var("SQL_PASSWORD").map_err(|_| Error::DatabaseCredentialsNotFound)?;
    Ok(DatabaseCredentials { username, password })
}
