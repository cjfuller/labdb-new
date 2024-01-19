use sqlx::FromRow;

#[derive(Clone, Debug, FromRow)]
#[allow(dead_code)]
pub struct Dino {
    pub id: i64,
    pub alias: Option<String>,
    pub description: Option<String>,
    pub entered_by: Option<String>,
}

impl super::SearchModel for Dino {
    fn table_name() -> &'static str {
        "dinos"
    }

    fn model_name() -> &'static str {
        "Dino"
    }

    fn selects() -> Vec<&'static str> {
        vec!["id", "alias", "entered_by", "description"]
    }

    fn id(&self) -> super::ModelID {
        super::ModelID {
            kind: Self::model_name(),
            id: self.id as i32,
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
        "description"
    }

    fn description(&self) -> &str {
        self.description.as_deref().unwrap_or_default()
    }

    fn sequence_field_name() -> Option<&'static str> {
        None
    }

    fn sequence(&self) -> Option<&str> {
        None
    }
}
