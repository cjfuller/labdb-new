#[derive(Clone, Debug, sqlx::FromRow)]
#[allow(dead_code)]
pub struct Sample {
    pub id: i32,
    pub entered_by: Option<String>,
    pub description: Option<String>,
    pub sample_alias: Option<String>,
}

impl super::SearchModel for Sample {
    fn table_name() -> &'static str {
        "samples"
    }

    fn model_name() -> &'static str {
        "Sample"
    }

    fn selects() -> Vec<&'static str> {
        vec!["id", "entered_by", "description", "sample_alias"]
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
        "sample_alias"
    }

    fn short_desc(&self) -> &str {
        self.sample_alias.as_deref().unwrap_or_default()
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
