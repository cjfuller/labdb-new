package env

import (
	"os"
)

var Dev = os.Getenv("DEV") == "1"
var Prod = !Dev
var SecretToken string
var SigningKey string
var DbURL = os.Getenv("DATABASE_URL")
var DebugDB = os.Getenv("DB_DEBUG") == "1"

func init() {
	if Dev {
		SigningKey = "development-key"
	} else {
		SigningKey = os.Getenv("SIGNING_KEY")
		if SigningKey == "" {
			panic("Must provide a signing key in prod.")
		}
	}
	if Dev {
		SecretToken = "development-token"
	} else {
		SecretToken = os.Getenv("SECRET_TOKEN")
		if SecretToken == "" {
			panic("Must provide a secret token in prod.")
		}
	}
}
