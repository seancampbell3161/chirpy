package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/seancampbell3161/chirpy/internal/database"
	"time"
)

func GenerateJWT(user database.User, secret string, expDuration time.Duration) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Second * expDuration)),
		Subject:   string(rune(user.ID)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signingSecret := []byte(secret)
	return token.SignedString(signingSecret)
}
