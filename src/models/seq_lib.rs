#[derive(Clone, Debug, sqlx::FromRow)]
#[allow(dead_code)]
pub struct SeqLib {
    pub id: i32,
    pub entered_by: Option<String>,
    pub alias: Option<String>,
    pub description: Option<String>,
    pub index_seq: Option<String>,
}

impl super::SearchModel for SeqLib {
    fn table_name() -> &'static str {
        "seq_libs"
    }

    fn model_name() -> &'static str {
        "SeqLib"
    }

    fn selects() -> Vec<&'static str> {
        vec!["id", "entered_by", "alias", "description", "index_seq"]
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
        "description"
    }

    fn description(&self) -> &str {
        self.description.as_deref().unwrap_or_default()
    }

    fn sequence_field_name() -> Option<&'static str> {
        Some("index_seq")
    }

    fn sequence(&self) -> Option<&str> {
        self.index_seq.as_deref()
    }
}
