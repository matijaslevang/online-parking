package main

import (
	"log"
	"net/http"

	"bank/internal/api"
	"bank/internal/auth"
	"bank/internal/config"
	"bank/internal/database"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}

	authenticator := auth.NewAuthenticator(cfg.KeycloakIssuerURL, cfg.KeycloakJWKSURL, cfg.KeycloakAudience)
	router := api.NewRouter(db, authenticator)

	log.Printf("bank service listening on :%s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
