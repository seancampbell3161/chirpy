package main

import (
	"encoding/json"
	"fmt"
	"github.com/seancampbell3161/chirpy/internal/auth"
	"net/http"
)

type tokenResponse struct {
	Token string `json:"token"`
}

func (cfg *apiConfig) refreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	user, err := cfg.DB.GetUserByID(1)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(401)
		return
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
