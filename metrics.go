package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) middlewareMetrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cfg.fileServerHits++
		next.ServeHTTP(w, req)
	})
}

func (cfg *apiConfig) getNumOfHitsHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Status-URI", "200")

	htmlTemplate := `
		<html>
		
		<body>
		<h1>Welcome, Chirpy Admin</h1>
		<p>Chirpy has been visited %d times!</p>
		</body>
		
		</html>
	`
	formattedTemplate := fmt.Sprintf(htmlTemplate, cfg.fileServerHits)
	_, err := w.Write([]byte(formattedTemplate))
	if err != nil {
		return
	}
}
