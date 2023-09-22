package main

import (
	"encoding/json"
	"github.com/seancampbell3161/chirpy/internal/auth"
	"log"
	"net/http"
	"strings"
)

type parameters struct {
	Body string `json:"body"`
}

type validResp struct {
	CleanedBody string `json:"cleaned_body"`
}

type errorResp struct {
	Error string `json:"error"`
}

func (cfg *apiConfig) newChirpHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.ValidateJWT(r, cfg.JwtSecret, "chirpy-access")
	if err != nil {
		w.WriteHeader(401)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
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

	if len(params.Body) > 140 {
		w.WriteHeader(400)
		badChirpResp := errorResp{
			Error: "Chirp is too long",
		}
		data, err := json.Marshal(badChirpResp)
		if err != nil {
			return
		}
		_, err = w.Write(data)
	} else {
		result := &params.Body
		result = filterBadWords(result)

		chirp, err := cfg.DB.CreateChirp(*result, userID)
		if err != nil {
			return
		}
		data, err := json.Marshal(chirp)

		w.WriteHeader(201)
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(data)
	}
}

func filterBadWords(body *string) *string {
	badWords := getBadWords()
	for _, word := range strings.Split(*body, " ") {
		if contains(badWords, strings.ToLower(word)) {
			*body = strings.Replace(*body, word, "****", -1)
		}
	}
	return body
}
