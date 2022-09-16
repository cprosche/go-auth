package utils

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
)

func ValidateToken(signedToken string) (int, error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&jwt.StandardClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		},
	)
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok {
		return 0, errors.New("couldn't parse claims")
	}
	if claims.ExpiresAt < time.Now().Unix() {
		return 0, errors.New("token expired")
	}
	id, err := strconv.Atoi(claims.Subject)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func GetRSAKey(key string) *rsa.PrivateKey {
	block, _ := pem.Decode([]byte(key))
	parseResult, _ := x509.ParsePKCS8PrivateKey(block.Bytes)
	rsaKey := parseResult.(*rsa.PrivateKey)
	return rsaKey
}
