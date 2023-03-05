use sqlx::postgres::PgPoolOptions;

use crate::env::Env;

pub fn database_url() -> String {
    if Env::dev() {
        "postgresql://postgres@localhost:5432/labdb".into()
    } else {
        let db_name = std::env::var("PG_DB").expect("Must provide the database name in prod.");
        let db_password =
            std::env::var("DB_PASSWORD").expect("Must provide a database password in prod.");
        format!(
            "postgresql://postgres@{}/{db_name}?password={}",
            urlencoding::encode("/cloudsql/labdb-io:northamerica-northeast1:labdb"),
            urlencoding::encode(&db_password),
        )
    }
}

pub async fn create_pool() -> sqlx::Result<sqlx::PgPool> {
    PgPoolOptions::new()
        .max_connections(50)
        .connect(&database_url())
        .await
}
