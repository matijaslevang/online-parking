package main

import (
	"log"
	"net/http"

	"psp/internal/api"
	"psp/internal/auth"
	"psp/internal/config"
	"psp/internal/database"
	"psp/internal/keycloak"
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

	if err := database.Migrate(db); err != nil {
		log.Fatalf("migrate database: %v", err)
	}

	authenticator := auth.NewAuthenticator(cfg.KeycloakIssuerURL, cfg.KeycloakJWKSURL, cfg.KeycloakAudience)
	keycloakAdmin := keycloak.NewAdminClient(cfg.KeycloakURL, cfg.KeycloakRealm, cfg.KeycloakClientID, cfg.KeycloakClientSecret)
	router := api.NewRouter(db, authenticator, keycloakAdmin)

	log.Printf("psp service listening on :%s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
