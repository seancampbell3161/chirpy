package main

import (
	"encoding/json"
	"fmt"
	"github.com/seancampbell3161/chirpy/internal/auth"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

type userParameters struct {
	Password         string `json:"password"`
	Email            string `json:"email"`
	ExpiresInSeconds int    `json:"expires_in_seconds,omitempty"`
}

type userLoginResponse struct {
	Email string `json:"email"`
	ID    int    `json:"id"`
	Token string `json:"token"`
}

func (cfg *apiConfig) loginUserHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	userParams := userParameters{}
	err := decoder.Decode(&userParams)
	if err != nil {
		fmt.Println(err)
	}

	user, err := cfg.DB.GetUserByEmail(userParams.Email)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(401)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userParams.Password))
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(401)
		return
	} else {
		token, err := auth.GenerateJWT(user, cfg.JwtSecret, time.Duration(userParams.ExpiresInSeconds))
		if err != nil {
			fmt.Println(err)
		}
		response := userLoginResponse{user.Email, user.ID, token}
		data, err := json.Marshal(response)
		if err != nil {
			fmt.Println(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		_, err = w.Write(data)
		if err != nil {
			fmt.Println(err)
		}
	}
}
