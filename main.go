package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"log"
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

func validateChirpHandler(writer http.ResponseWriter, request *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type validResp struct {
		Valid bool `json:"valid"`
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
		respBody := validResp{
			Valid: true,
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

func main() {
	//mux := http.NewServeMux()
	myConfig := apiConfig{fileServerHits: 0}
	r := chi.NewRouter()

	//r.Handle("/app/", myConfig.middlewareMetrics(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	fsHandler := myConfig.middlewareMetrics(http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	r.Handle("/app", fsHandler)
	r.Handle("/app/*", fsHandler)

	//r.HandleFunc("/healthz", statusHandler)
	//r.HandleFunc("/metrics", myConfig.getNumOfHitsHandler)

	apiRouter := chi.NewRouter()
	apiRouter.Get("/healthz", statusHandler)
	apiRouter.Post("/validate_chirp", validateChirpHandler)

	adminRouter := chi.NewRouter()
	adminRouter.Get("/metrics", myConfig.getNumOfHitsHandler)

	r.Mount("/api/", apiRouter)
	r.Mount("/admin", adminRouter)

	corsMux := middlewareCORS(r)

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
