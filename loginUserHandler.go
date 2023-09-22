package main

import (
	"encoding/json"
	"fmt"
	"github.com/seancampbell3161/chirpy/internal/auth"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type userParameters struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type userLoginResponse struct {
	Email        string `json:"email"`
	ID           int    `json:"id"`
	IsChirpyRed  bool   `json:"is_chirpy_red"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
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
		accessToken, err := auth.GenerateJWT(user, cfg.JwtSecret, cfg.AccessExp, "chirpy-access")
		if err != nil {
			fmt.Println(err)
		}
		refreshToken, err := auth.GenerateJWT(user, cfg.JwtSecret, cfg.RefreshExp, "chirpy-refresh")
		if err != nil {
			fmt.Println(err)
		}
		response := userLoginResponse{
			user.Email,
			user.ID,
			user.IsChirpyRed,
			accessToken,
			refreshToken,
		}
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
