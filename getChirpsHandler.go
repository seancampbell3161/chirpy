package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
)

func (cfg *apiConfig) getChirpsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	authID := r.URL.Query().Get("author_id")
	sortParam := r.URL.Query().Get("sort")

	chirps, err := cfg.DB.GetChirps(&authID)
	if err != nil {
		w.WriteHeader(500)
		fmt.Print(err)
	}

	if sortParam == "asc" {
		sort.SliceStable(chirps, func(i, j int) bool {
			return chirps[i].ID < chirps[j].ID
		})
	} else if sortParam == "desc" {
		sort.SliceStable(chirps, func(i, j int) bool {
			return chirps[i].ID > chirps[j].ID
		})
	}
	w.WriteHeader(200)
	data, err := json.Marshal(chirps)
	_, err = w.Write(data)
	if err != nil {
		return
	}
}
