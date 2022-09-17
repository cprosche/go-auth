package inits

import (
	"crypto/rsa"

	"github.com/joho/godotenv"
)

// define globals
var (
	PORT        string
	PRIVATE_KEY string
	AUTH_DSN    string
	CURRENT_ENV string
	RSA_KEY     *rsa.PrivateKey
)

func LoadEnv() {
	godotenv.Load(".env")
}
