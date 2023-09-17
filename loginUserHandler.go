package main

import (
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type userParameters struct {
	Password           string `json:"password"`
	Email              string `json:"email"`
	Expires_in_seconds int    `json:"expires_in_seconds"`
}

func (cfg *apiConfig) loginUserHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	userParams := userParameters{}
	err := decoder.Decode(&userParams)
	if err != nil {
		w.WriteHeader(500)
	}

	user, err := cfg.DB.GetUserByEmail(userParams.Email)
	if err != nil {
		w.WriteHeader(401)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userParams.Password))
	if err != nil {
		w.WriteHeader(401)
	} else {
		response := userResponse{user.Email, user.ID}
		data, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(500)
		}
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(data)
		if err != nil {
			return
		}
	}
}
