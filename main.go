package main

import (
	"flag"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/seancampbell3161/chirpy/internal/database"
	"log"
	"net/http"
	"os"
	"time"
)

type apiConfig struct {
	fileServerHits int
	DB             *database.DB
	JwtSecret      string
	AccessExp      time.Duration
	RefreshExp     time.Duration
}

func main() {
	dbg := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()
	if *dbg == true {
		err := os.Remove("database.json")
		if err != nil {
			log.Fatal(err)
		}
	}
	db, err := database.NewDB("database.json")
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}

	err = godotenv.Load()
	if err != nil {
		return
	}
	jwtSecret := os.Getenv("JWT_SECRET")

	appConfig := apiConfig{
		fileServerHits: 0,
		DB:             db,
		JwtSecret:      jwtSecret,
		AccessExp:      time.Hour,
		RefreshExp:     time.Hour * 24 * 60,
	}
	r := chi.NewRouter()

	fsHandler := appConfig.middlewareMetrics(http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	r.Handle("/app", fsHandler)
	r.Handle("/app/*", fsHandler)

	apiRouter := chi.NewRouter()
	apiRouter.Get("/healthz", statusHandler)
	apiRouter.Get("/chirps", appConfig.getChirpsHandler)
	apiRouter.Get("/chirps/{chirpID}", appConfig.getChirpHandler)
	apiRouter.Post("/chirps", appConfig.newChirpHandler)
	apiRouter.Delete("/chirps/{chirpID}", appConfig.deleteChirpHandler)
	apiRouter.Post("/users", appConfig.createUserHandler)
	apiRouter.Put("/users", appConfig.updateUserHandler)
	apiRouter.Post("/login", appConfig.loginUserHandler)
	apiRouter.Post("/refresh", appConfig.refreshTokenHandler)
	apiRouter.Post("/revoke", appConfig.revokeRefreshTokenHandler)

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
