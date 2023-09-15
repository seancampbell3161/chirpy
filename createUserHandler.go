package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type userParams struct {
	Email string `json:"email"`
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

	user, err := cfg.DB.CreateUser(userParameters.Email)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	data, err := json.Marshal(user)
	_, err = w.Write(data)
	if err != nil {
		w.WriteHeader(500)
		return
	}
}
