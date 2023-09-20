package main

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/seancampbell3161/chirpy/internal/auth"
	"net/http"
	"strconv"
	"strings"
)

type tokenResponse struct {
	Token string `json:"token"`
}

type userClaims struct {
	Password string
	Email    string
	ID       int
	jwt.RegisteredClaims
}

func (cfg *apiConfig) refreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	// refactor this bc you're doing this in the auth func
	tokenString := strings.Split(r.Header.Get("Authorization"), "Bearer ")[1]
	token, err := jwt.ParseWithClaims(tokenString, &userClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JwtSecret), nil
	})
	if err != nil {
		fmt.Println("parse w claims", err)
		return
	}

	userID, err := token.Claims.GetSubject()
	if err != nil {
		fmt.Println(err)
	}

	userIDint, err := strconv.Atoi(userID)

	user, err := cfg.DB.GetUserByID(userIDint)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(401)
		return
	}

	revokedToken, err := cfg.DB.GetRevokedRefreshToken(tokenString)
	if err != nil {
		fmt.Println(err)
	}
	if len(revokedToken) > 0 && revokedToken == tokenString {
		fmt.Println("token has been revoked")
		w.WriteHeader(401)
	}

	accessToken, err := auth.RefreshAccessJWT(r, cfg.JwtSecret, user, cfg.AccessExp)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(401)
		return
	}

	response := tokenResponse{Token: accessToken}
	data, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	_, err = w.Write(data)
	if err != nil {
		fmt.Println(err)
	}
}
