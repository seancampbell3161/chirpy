package main

import (
	"fmt"
	"net/http"
)

type apiConfig struct {
	fileServerHits int
}

func middlewareCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if req.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, req)
	})
}

func statusHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Status-URI", "200")
	_, err := w.Write([]byte("OK"))
	if err != nil {
		return
	}
}

func (cfg *apiConfig) middlewareMetrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cfg.fileServerHits++
		next.ServeHTTP(w, req)
	})
}

func (cfg *apiConfig) getNumOfHitsHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Status-URI", "200")

	numOfHits := fmt.Sprintf("Hits: %d", cfg.fileServerHits)
	_, err := w.Write([]byte(numOfHits))
	if err != nil {
		return
	}
}

func main() {
	mux := http.NewServeMux()
	myConfig := apiConfig{fileServerHits: 0}

	mux.Handle("/app/", myConfig.middlewareMetrics(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))

	mux.HandleFunc("/healthz", statusHandler)
	mux.HandleFunc("/metrics", myConfig.getNumOfHitsHandler)

	corsMux := middlewareCORS(mux)

	server := &http.Server{
		Handler: corsMux,
		Addr:    ":8080",
	}
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
		return
	}
}
