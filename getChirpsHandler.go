package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (cfg *apiConfig) getChirpsHandler(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	chirps, err := cfg.DB.GetChirps()
	if err != nil {
		writer.WriteHeader(500)
		fmt.Print(err)
	}
	writer.WriteHeader(200)
	data, err := json.Marshal(chirps)
	_, err = writer.Write(data)
	if err != nil {
		return
	}
}
