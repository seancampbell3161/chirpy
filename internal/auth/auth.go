package auth

import (
	"errors"
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

func GenerateJWT(user database.User, secret string, expDuration time.Duration, tokenIssuer string) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:    tokenIssuer,
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expDuration)),
		Subject:   strconv.FormatInt(int64(user.ID), 10),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signingSecret := []byte(secret)
	return token.SignedString(signingSecret)
}

func RefreshAccessJWT(r *http.Request, secret string, user database.User, exp time.Duration) (string, error) {
	tokenString := r.Header.Get("Authorization")
	tokenString = strings.Split(tokenString, "Bearer ")[1]

	token, err := jwt.ParseWithClaims(tokenString, &userClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		fmt.Println("parse w claims: ", err)
		return "", err
	}

	if iss, err := token.Claims.GetIssuer(); iss != "chirpy-refresh" {
		if err != nil {
			fmt.Println("error getting issuer")
			return "", errors.New("error getting issuer")
		}
		fmt.Println("issuer is not valid for refresh")
		return "", errors.New("issuer not valid for refresh")
	}

	return GenerateJWT(user, secret, exp, "chirpy-access")
}

func RevokeRefreshJWT(r *http.Request, secret string) (string, error) {
	tokenString := strings.Split(r.Header.Get("Authorization"), "Bearer ")[1]
	token, err := jwt.ParseWithClaims(tokenString, &userClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		fmt.Println("parse w claims", err)
		return "", err
	}
	if iss, err := token.Claims.GetIssuer(); iss != "chirpy-refresh" {
		if err != nil {
			fmt.Println("issuer is snot valid for revoke")
			return "", errors.New("issuer not valid for revoke")
		}
	}
	return token.Raw, nil
}

func ValidateJWT(r *http.Request, secret string, tokenIssuer string) (int, error) {
	tokenString := r.Header.Get("Authorization")
	tokenString = strings.Split(tokenString, "Bearer ")[1]

	token, err := jwt.ParseWithClaims(tokenString, &userClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		fmt.Println("Get token issuer: ", err)
	}
	if issuer != tokenIssuer {
		return 0, errors.New("invalid issuer")
	}

	stringID, err := token.Claims.GetSubject()
	userID, err := strconv.Atoi(stringID)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return userID, nil
}
