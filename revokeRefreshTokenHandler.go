package main

import (
	"fmt"
	"github.com/seancampbell3161/chirpy/internal/auth"
	"net/http"
)

func (cfg *apiConfig) revokeRefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	tokenString, err := auth.RevokeRefreshJWT(r, cfg.JwtSecret)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(401)
		return
	}

	fmt.Println(tokenString)
	err = cfg.DB.AddRevokedRefreshToken(tokenString)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(200)
}
