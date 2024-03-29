package main

import (
	"encoding/json"
	"fmt"
	"github.com/seancampbell3161/chirpy/internal/auth"
	"log"
	"net/http"
)

type userParams struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type userResponse struct {
	Email       string `json:"email"`
	ID          int    `json:"id"`
	IsChirpyRed bool   `json:"is_chirpy_red"`
}

func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	userParameters := userParams{}
	err := decoder.Decode(&userParameters)
	if err != nil {
		w.WriteHeader(500)
		respBody := errorResp{
			Error: "Something went wrong",
		}
		data, err := json.Marshal(respBody)
		if err != nil {
			log.Printf("Error marshalling response: %s", err)
		}
		_, err = w.Write(data)
		return
	}
	userParameters.Password = auth.GenerateHashedPassword(&userParameters.Password)

	userResult, err := cfg.DB.CreateUser(userParameters.Email, userParameters.Password)
	response := userResponse{userResult.Email, userResult.ID, userResult.IsChirpyRed}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	data, err := json.Marshal(response)
	_, err = w.Write(data)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}
}
