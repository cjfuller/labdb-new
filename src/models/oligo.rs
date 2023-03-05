#[derive(Clone, Debug, sqlx::FromRow)]
#[allow(dead_code)]
pub struct Oligo {
    pub id: i32,
    pub entered_by: Option<String>,
    pub oligoalias: Option<String>,
    pub purpose: Option<String>,
    pub sequence: Option<String>,
}

impl super::SearchModel for Oligo {
    fn table_name() -> &'static str {
        "oligos"
    }

    fn model_name() -> &'static str {
        "Oligo"
    }

    fn selects() -> Vec<&'static str> {
        vec!["id", "entered_by", "oligoalias", "purpose", "sequence"]
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
        "oligoalias"
    }

    fn short_desc(&self) -> &str {
        &self.oligoalias.as_deref().unwrap_or_default()
    }

    fn description_field_name() -> &'static str {
        "purpose"
    }

    fn description(&self) -> &str {
        &self.purpose.as_deref().unwrap_or_default()
    }

    fn sequence_field_name() -> Option<&'static str> {
        Some("sequence")
    }

    fn sequence(&self) -> Option<&str> {
        self.sequence.as_deref()
    }
}
