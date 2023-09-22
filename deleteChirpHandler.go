package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/seancampbell3161/chirpy/internal/auth"
	"net/http"
	"strconv"
)

func (cfg *apiConfig) deleteChirpHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.ValidateJWT(r, cfg.JwtSecret, "chirpy-access")
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(401)
		return
	}

	chirpID := chi.URLParam(r, "chirpID")
	chirpIDString, err := strconv.Atoi(chirpID)
	if err != nil {
		fmt.Println("error converting string")
		w.WriteHeader(500)
		return
	}

	err = cfg.DB.DeleteChirp(chirpIDString, userID)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(403)
		return
	}

	w.WriteHeader(200)
}
