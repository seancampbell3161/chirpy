package main

import (
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

type userParams struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type userResponse struct {
	Email string `json:"email"`
	ID    int    `json:"id"`
	Token string `json:"token"`
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
	userParameters.Password = generateHashedPassword(&userParameters.Password)

	userResult, err := cfg.DB.CreateUser(userParameters.Email, userParameters.Password)
	response := userResponse{userResult.Email, userResult.ID, ""}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	data, err := json.Marshal(response)
	_, err = w.Write(data)
	if err != nil {
		w.WriteHeader(500)
		return
	}
}

func generateHashedPassword(password *string) string {
	data, err := bcrypt.GenerateFromPassword([]byte(*password), 12)
	if err != nil {
		log.Fatal(err)
	}
	return string(data)
}
