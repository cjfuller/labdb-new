#[derive(Clone, Debug, sqlx::FromRow)]
#[allow(dead_code)]
pub struct RNAiClone {
    id: i32,
    alias: Option<String>,
    description: Option<String>,
    entered_by: Option<String>,
}

impl super::SearchModel for RNAiClone {
    fn table_name() -> &'static str {
        "rnai_clones"
    }

    fn model_name() -> &'static str {
        "RNAiClone"
    }

    fn selects() -> Vec<&'static str> {
        vec!["id", "alias", "description", "entered_by"]
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
        &self.alias.as_deref().unwrap_or_default()
    }

    fn description_field_name() -> &'static str {
        "description"
    }

    fn description(&self) -> &str {
        &self.description.as_deref().unwrap_or_default()
    }

    fn sequence_field_name() -> Option<&'static str> {
        None
    }

    fn sequence(&self) -> Option<&str> {
        None
    }
}
