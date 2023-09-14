package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func contains(slice []string, word string) bool {
	for _, item := range slice {
		if item == word {
			return true
		}
	}
	return false
}

func getBadWords() []string {
	return []string{
		"kerfuffle",
		"sharbert",
		"fornax",
	}
}

func validateChirpHandler(writer http.ResponseWriter, request *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type validResp struct {
		Cleaned_body string `json:"cleaned_body"`
	}

	type errorResp struct {
		Error string `json:"error"`
	}

	decoder := json.NewDecoder(request.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		writer.WriteHeader(500)
		respBody := errorResp{
			Error: "Something went wrong",
		}
		data, err := json.Marshal(respBody)
		if err != nil {
			log.Printf("Error marshalling response: %s", err)
		}
		_, err = writer.Write(data)
		return
	}

	if len(params.Body) > 140 {
		writer.WriteHeader(400)
		badChirpResp := errorResp{
			Error: "Chirp is too long",
		}
		data, err := json.Marshal(badChirpResp)
		if err != nil {
			return
		}
		_, err = writer.Write(data)
	} else {
		result := &params.Body
		badWords := getBadWords()
		for _, word := range strings.Split(params.Body, " ") {
			if contains(badWords, strings.ToLower(word)) {
				*result = strings.Replace(*result, word, "****", -1)
			}
		}
		respBody := validResp{
			Cleaned_body: *result,
		}
		data, err := json.Marshal(respBody)
		if err != nil {
			log.Printf("Error marshalling response: %s", err)
			writer.WriteHeader(500)
			return
		}
		writer.WriteHeader(200)
		writer.Header().Set("Content-Type", "application/json")
		_, err = writer.Write(data)
	}
}
