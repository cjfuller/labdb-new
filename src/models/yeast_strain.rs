#[derive(Clone, Debug, sqlx::FromRow)]
#[allow(dead_code)]
pub struct YeastStrain {
    pub id: i32,
    pub entered_by: Option<String>,
    pub comments: Option<String>,
    pub strainalias: Option<String>,
    pub sequence: Option<String>,
}

impl super::SearchModel for YeastStrain {
    fn table_name() -> &'static str {
        "yeaststrains"
    }

    fn model_name() -> &'static str {
        // yes, this is correct :(
        "Yeaststrain"
    }

    fn selects() -> Vec<&'static str> {
        vec!["id", "entered_by", "strainalias", "comments", "sequence"]
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
        "strainalias"
    }

    fn short_desc(&self) -> &str {
        self.strainalias.as_deref().unwrap_or_default()
    }

    fn description_field_name() -> &'static str {
        "comments"
    }

    fn description(&self) -> &str {
        self.comments.as_deref().unwrap_or_default()
    }

    fn sequence_field_name() -> Option<&'static str> {
        Some("sequence")
    }

    fn sequence(&self) -> Option<&str> {
        self.sequence.as_deref()
    }
}
