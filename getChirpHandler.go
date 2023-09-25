package main

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func (cfg *apiConfig) getChirpHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := chi.URLParam(r, "chirpID")

	authID := ""
	if num, err := strconv.Atoi(id); err == nil {
		chirps, err := cfg.DB.GetChirps(&authID)
		if err != nil {
			w.WriteHeader(500)
		}
		if num > len(chirps) {
			w.WriteHeader(404)
		} else {
			chirp := chirps[num-1]
			data, err := json.Marshal(chirp)
			w.WriteHeader(200)
			_, err = w.Write(data)
			if err != nil {
				return
			}
		}
	}
}
