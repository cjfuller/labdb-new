use std::env;

use once_cell::sync::Lazy;

#[derive(Clone, Copy, Debug, Eq, PartialEq)]
pub enum Env {
    Dev,
    Prod,
}

impl Env {
    pub fn dev() -> bool {
        Self::default() == Self::Dev
    }
    pub fn prod() -> bool {
        Self::default() == Self::Prod
    }
}

impl Default for Env {
    fn default() -> Self {
        if &env::var("DEV").unwrap_or_else(|_| "0".into()) == "1" {
            Env::Dev
        } else {
            Env::Prod
        }
    }
}

pub static SIGNING_KEY: Lazy<String> = Lazy::new(|| {
    if Env::dev() {
        "development-key".into()
    } else {
        env::var("SIGNING_KEY").expect("Must provide a signing key in prod.")
    }
});

pub static SECRET_TOKEN: Lazy<String> = Lazy::new(|| {
    if Env::dev() {
        "development-token-00000000000000000000000000000000000000000000000".into()
    } else {
        env::var("SECRET_TOKEN").expect("Must provide a secret token in prod.")
    }
});
