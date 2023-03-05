#!/usr/bin/env nu

def current_hash [] {
    git log -1 --format=%h | str trim
}

def build_container [] {
    let sha = current_hash
    (gcloud builds submit --region northamerica-northeast1
        --config cloudbuild.yml $"--substitutions=SHORT_SHA=($sha)")
}

def main [] {
    build_container
    echo "Now update the cloud run deployment in the GCP console to point to this new container."
}
