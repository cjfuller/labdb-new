#[derive(Clone, Debug, sqlx::FromRow)]
#[allow(dead_code)]
pub struct CellLine {
    pub id: i32,
    pub entered_by: Option<String>,
    pub description: Option<String>,
    pub line_alias: Option<String>,
    pub sequence: Option<String>,
}

impl super::SearchModel for CellLine {
    fn table_name() -> &'static str {
        "lines"
    }

    fn model_name() -> &'static str {
        "Line"
    }

    fn selects() -> Vec<&'static str> {
        vec!["id", "entered_by", "description", "line_alias", "sequence"]
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
        "line_alias"
    }

    fn short_desc(&self) -> &str {
        &self.line_alias.as_deref().unwrap_or_default()
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
