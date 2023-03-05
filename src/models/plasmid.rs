use sqlx::FromRow;

use super::{ModelID, SearchModel};

#[derive(Clone, Debug, FromRow)]
#[allow(dead_code)]
pub struct Plasmid {
    pub id: i32,
    pub alias: Option<String>,
    pub description: Option<String>,
    pub sequence: Option<String>,
    pub creator: Option<String>,
}

impl SearchModel for Plasmid {
    fn table_name() -> &'static str {
        "plasmids"
    }

    fn model_name() -> &'static str {
        "Plasmid"
    }

    fn selects() -> Vec<&'static str> {
        vec!["id", "alias", "description", "sequence", "creator"]
    }

    fn id(&self) -> ModelID {
        ModelID {
            kind: Self::model_name(),
            id: self.id,
        }
    }

    fn owner_field_name() -> &'static str {
        "creator"
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
        Some("sequence")
    }

    fn sequence(&self) -> Option<&str> {
        self.sequence.as_deref()
    }
}
