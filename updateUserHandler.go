package main

import (
	"encoding/json"
	"fmt"
	"github.com/seancampbell3161/chirpy/internal/auth"
	"net/http"
)

func (cfg *apiConfig) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.ValidateJWT(r, cfg.JwtSecret)
	if err != nil {
		w.WriteHeader(401)
		return
	}

	userParams := userParameters{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&userParams)
	if err != nil {
		fmt.Println(err)
		return
	}

	userParams.Password = auth.GenerateHashedPassword(&userParams.Password)

	user, err := cfg.DB.UpdateUser(userID, userParams.Email, userParams.Password)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	response := userResponse{ID: user.ID, Email: user.Email}
	data, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	_, err = w.Write(data)
	if err != nil {
		fmt.Println(err)
		return
	}
}
