package inits

import (
	"crypto/rsa"

	env "github.com/joho/godotenv"
)

var PORT string
var PRIVATE_KEY string
var AUTH_DSN string
var CURRENT_ENV string
var HMAC_KEY string
var RSA_KEY *rsa.PrivateKey

func LoadEnv() {
	err := env.Load(".env")
	if err != nil {
		panic("Error loading environment variables")
	}
}
