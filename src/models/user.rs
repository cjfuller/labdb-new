#[derive(Clone, Debug, sqlx::FromRow)]
#[allow(dead_code)]
pub struct User {
    id: i32,
    email: Option<String>,
    name: Option<String>,
    notes: Option<String>,
}

impl super::SearchModel for User {
    fn table_name() -> &'static str {
        "users"
    }

    fn model_name() -> &'static str {
        "User"
    }

    fn selects() -> Vec<&'static str> {
        vec!["id", "email", "name", "notes"]
    }

    fn id(&self) -> super::ModelID {
        super::ModelID {
            kind: Self::model_name(),
            id: self.id,
        }
    }

    fn owner_field_name() -> &'static str {
        "name"
    }

    fn short_desc_field_name() -> &'static str {
        "email"
    }

    fn short_desc(&self) -> &str {
        self.email.as_deref().unwrap_or_default()
    }

    fn description_field_name() -> &'static str {
        "notes"
    }

    fn description(&self) -> &str {
        self.notes.as_deref().unwrap_or_default()
    }

    fn sequence_field_name() -> Option<&'static str> {
        None
    }

    fn sequence(&self) -> Option<&str> {
        None
    }
}
