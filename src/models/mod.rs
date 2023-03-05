use regex::Regex;
use serde::Serialize;
use sqlx::{Database, FromRow, PgConnection, Postgres};

use crate::errors::Result;

pub mod antibody;
pub mod bacterial_strain;
pub mod cell_line;
pub mod oligo;
pub mod plasmid;
pub mod rnai_clone;
pub mod sample;
pub mod seq_lib;
pub mod user;
pub mod yeast_strain;

#[derive(Clone, Copy, Debug, Serialize)]
pub struct ModelID {
    kind: &'static str,
    id: i32,
}

pub trait Model {
    fn created_at(&self) -> chrono::NaiveDateTime;
    fn updated_at(&self) -> chrono::NaiveDateTime;
}

// TODO: write a macro to generate the projection for search from the base model.
pub trait SearchModel {
    fn table_name() -> &'static str;
    fn model_name() -> &'static str;
    fn selects() -> Vec<&'static str>;

    fn id(&self) -> ModelID;

    fn owner_field_name() -> &'static str;

    fn short_desc_field_name() -> &'static str;
    fn short_desc(&self) -> &str;

    fn description_field_name() -> &'static str;
    fn description(&self) -> &str;

    fn sequence_field_name() -> Option<&'static str>;
    fn sequence(&self) -> Option<&str>;
}

pub type QueryResult = (&'static str, i32);

impl From<ModelID> for QueryResult {
    fn from(value: ModelID) -> Self {
        (value.kind, value.id)
    }
}

pub async fn owner_search_query<M>(conn: &mut PgConnection, name: &str) -> Result<Vec<QueryResult>>
where
    for<'r> M: SearchModel + Send + Unpin + FromRow<'r, <Postgres as Database>::Row>,
{
    Ok(sqlx::query_as::<_, M>(&format!(
        "SELECT {} FROM {} WHERE {} = $1 ORDER BY id DESC",
        M::selects().join(","),
        M::table_name(),
        M::owner_field_name(),
    ))
    .bind(name)
    .fetch_all(conn)
    .await?
    .into_iter()
    .map(|it| it.id().into())
    .collect())
}

pub async fn search_query<M>(
    conn: &mut PgConnection,
    re: &Regex,
    include_sequence: bool,
) -> Result<Vec<QueryResult>>
where
    for<'r> M: SearchModel + Send + Unpin + FromRow<'r, <Postgres as Database>::Row>,
{
    let select_parts = M::selects().join(",");
    // Note: the strings being interpolated here ultimately come from code, not user input (note the
    // &'static str type), so ok to insert these directly in the query.
    let items = sqlx::query_as::<_, M>(&format!(
        "SELECT {select_parts} FROM {} ORDER BY id DESC",
        M::table_name()
    ))
    .fetch_all(conn)
    .await?;

    Ok(items
        .into_iter()
        .filter(|it| {
            re.is_match(it.short_desc())
                || re.is_match(it.description())
                || (include_sequence && re.is_match(it.sequence().unwrap_or_default()))
        })
        .map(|it| it.id().into())
        .collect())
}

#[macro_export]
macro_rules! model_eval {
    ($model_name:expr, $typevar:ident, $code:block) => {{
        match $model_name.to_ascii_lowercase().as_ref() {
            "plasmid" | "plasmids" => {
                type $typevar = $crate::models::plasmid::Plasmid;
                $code
            }
            "oligo" | "oligos" => {
                type $typevar = $crate::models::oligo::Oligo;
                $code
            }
            "line" | "lines" => {
                type $typevar = $crate::models::cell_line::CellLine;
                $code
            }
            "sample" | "samples" => {
                type $typevar = $crate::models::sample::Sample;
                $code
            }
            "bacterium" | "bacteria" => {
                type $typevar = $crate::models::bacterial_strain::BacterialStrain;
                $code
            }
            "yeaststrain" | "yeaststrains" => {
                type $typevar = $crate::models::yeast_strain::YeastStrain;
                $code
            }
            "user" | "users" => {
                type $typevar = $crate::models::user::User;
                $code
            }
            "antibody" | "antibodies" => {
                type $typevar = $crate::models::antibody::Antibody;
                $code
            }
            "rnaiclone" | "rnaiclonse" | "rnai_clone" | "rnai_clones" => {
                type $typevar = $crate::models::rnai_clone::RNAiClone;
                $code
            }
            "seqlib" | "seqlibs" | "seq_lib" | "seq_libs" => {
                type $typevar = $crate::models::seq_lib::SeqLib;
                $code
            }
            _ => Err($crate::errors::Error::str(format!(
                "Unknown model type: {}",
                $model_name
            ))),
        }
    }};
}

// pub trait Entity: Model {
//     fn number(&self) -> u64;
//     fn kind(&self) -> &str;
//     fn name(&self) -> Option<&str>;
//     fn short_desc(&self) -> Option<&str>;
//     fn desc(&self) -> Option<&str>;
//     fn sequence(&self) -> Option<&str>;
//     fn owner_field_name(&self) -> Option<&str>;
//     fn core_links(&self) -> Option<CoreLinks>;
//     fn core_info_sections(&self) -> Vec<InfoSection>;
//     fn sequence_info(&self) -> Option<SequenceInfo>;
//     fn supplemental_fields(&self) -> Vec<FieldDef>;
//     fn as_resource_def(&self) -> String;
// }

// #[derive(Debug, Serialize)]
// pub struct FieldDef {
//     name: String,
//     lookup: String,
//     r#type: String,
// }

// #[derive(Debug, Serialize)]
// pub struct InfoSection {
//     name: String,
//     preformatted: bool,
//     lookup: String,
//     single: bool,
//     fields: Vec<FieldDef>,
//     inline_value: Option<String>,
// }

// #[derive(Debug, Serialize)]
// pub struct SequenceInfo {
//     sequence: FieldDef,
//     verified: FieldDef,
// }

// #[derive(Debug, Serialize)]
// pub struct CoreLinks {
//     lookup: String,
//     name: String,
//     links: Vec<String>,
// }

// #[derive(Debug, Serialize)]
// pub struct CoreInfoSections {}

// #[derive(Debug, Serialize)]
// pub struct ResourceDef {}

#[derive(Clone, Debug, Serialize, sqlx::FromRow)]
pub struct User {
    pub id: i32,
    // TODO(colin): migrate these to be non-null.
    pub email: Option<String>,
    pub name: Option<String>,
    pub auth_read: Option<bool>,
    pub auth_write: Option<bool>,
    pub auth_admin: Option<bool>,
    pub notes: Option<String>,
    pub created_at: chrono::NaiveDateTime,
    pub updated_at: chrono::NaiveDateTime,
}

impl User {
    pub async fn find_by_email(conn: &mut PgConnection, email: &str) -> Option<User> {
        let resp = sqlx::query_as::<_, User>("SELECT * FROM users WHERE email = $1")
            .bind(email)
            .fetch_optional(conn)
            .await;

        resp.ok().flatten()
    }
}
