use sqlx::FromRow;

#[derive(Clone, Debug, FromRow)]
#[allow(dead_code)]
pub struct Antibody {
    pub id: i32,
    pub entered_by: Option<String>,
    pub alias: Option<String>,
    pub comments: Option<String>,
}

impl super::SearchModel for Antibody {
    fn table_name() -> &'static str {
        "antibodies"
    }

    fn model_name() -> &'static str {
        "Antibody"
    }

    fn selects() -> Vec<&'static str> {
        vec!["id", "entered_by", "alias", "comments"]
    }

    fn id(&self) -> super::ModelID {
        super::ModelID {
            kind: Self::model_name(),
            id: self.id,
        }
    }

    fn owner_field_name() -> &'static str {
        "entered_by"
    }

    fn short_desc_field_name() -> &'static str {
        "alias"
    }

    fn short_desc(&self) -> &str {
        self.alias.as_deref().unwrap_or_default()
    }

    fn description_field_name() -> &'static str {
        "comments"
    }

    fn description(&self) -> &str {
        &self.comments.as_deref().unwrap_or_default()
    }

    fn sequence_field_name() -> Option<&'static str> {
        None
    }

    fn sequence(&self) -> Option<&str> {
        None
    }
}
