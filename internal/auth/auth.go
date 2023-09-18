package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/seancampbell3161/chirpy/internal/database"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type userClaims struct {
	Password string
	Email    string
	ID       int
	jwt.RegisteredClaims
}

func GenerateHashedPassword(password *string) string {
	data, err := bcrypt.GenerateFromPassword([]byte(*password), 12)
	if err != nil {
		log.Fatal(err)
	}
	return string(data)
}

func GenerateJWT(user database.User, secret string, expDuration time.Duration) (string, error) {
	if expDuration == time.Duration(0) {
		expDuration = time.Hour * 24
	}
	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expDuration)),
		Subject:   strconv.FormatInt(int64(user.ID), 10),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signingSecret := []byte(secret)
	return token.SignedString(signingSecret)
}

func ValidateJWT(r *http.Request, secret string) (int, error) {
	tokenString := r.Header.Get("Authorization")
	tokenString = strings.Split(tokenString, "Bearer ")[1]

	token, err := jwt.ParseWithClaims(tokenString, &userClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	stringID, err := token.Claims.GetSubject()
	userID, err := strconv.Atoi(stringID)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return userID, nil
}
