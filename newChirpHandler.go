package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type parameters struct {
	Body string `json:"body"`
}

type validResp struct {
	Cleaned_body string `json:"cleaned_body"`
}

type errorResp struct {
	Error string `json:"error"`
}

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

func (cfg *apiConfig) newChirpHandler(writer http.ResponseWriter, request *http.Request) {
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
		result = filterBadWords(result)

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

func filterBadWords(body *string) *string {
	badWords := getBadWords()
	for _, word := range strings.Split(*body, " ") {
		if contains(badWords, strings.ToLower(word)) {
			*body = strings.Replace(*body, word, "****", -1)
		}
	}
	return body
}
