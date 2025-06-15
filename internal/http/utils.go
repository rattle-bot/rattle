package http

import (
	"crypto/rand"
	"crypto/rsa"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ilyxenc/rattle/internal/logger"
)

var (
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
)

func GenerateAccessToken(id string) (string, error) {
	claims := jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(time.Minute * 30).Unix(), // 30 minutes for access_token
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	t, err := token.SignedString(PrivateKey)
	if err != nil {
		return "", err
	}

	return t, nil
}

func InitKeys() {
	// Generation RSA keys
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		logger.Log.Panicf("Error generating private key: %v", err)
	}

	// Setting keys as global
	PrivateKey = privateKey
	PublicKey = &privateKey.PublicKey
}
