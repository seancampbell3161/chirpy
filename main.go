package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/seancampbell3161/chirpy/internal/database"
	"log"
	"net/http"
)

type apiConfig struct {
	fileServerHits int
	DB             *database.DB
}

func main() {
	db, err := database.NewDB("database.json")
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	appConfig := apiConfig{
		fileServerHits: 0,
		DB:             db,
	}
	r := chi.NewRouter()

	fsHandler := appConfig.middlewareMetrics(http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	r.Handle("/app", fsHandler)
	r.Handle("/app/*", fsHandler)

	apiRouter := chi.NewRouter()
	apiRouter.Get("/healthz", statusHandler)
	apiRouter.Get("/chirps", appConfig.getChirpsHandler)
	apiRouter.Post("/chirps", appConfig.newChirpHandler)

	adminRouter := chi.NewRouter()
	adminRouter.Get("/metrics", appConfig.getNumOfHitsHandler)

	r.Mount("/api/", apiRouter)
	r.Mount("/admin", adminRouter)

	corsMux := middlewareCORS(r)

	server := &http.Server{
		Handler: corsMux,
		Addr:    ":8080",
	}
	err = server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
		return
	}
}
