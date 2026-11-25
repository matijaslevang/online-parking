package config

import (
	"errors"
	"os"
)

type Config struct {
	Port              string
	DatabaseURL       string
	KeycloakIssuerURL string
	KeycloakJWKSURL   string
	KeycloakAudience  string

	// Used to create merchant users through the Keycloak Admin API.
	KeycloakURL          string
	KeycloakRealm        string
	KeycloakClientID     string
	KeycloakClientSecret string
}

func Load() (Config, error) {
	cfg := Config{
		Port:              getEnv("PORT", "8080"),
		DatabaseURL:       os.Getenv("DATABASE_URL"),
		KeycloakIssuerURL: os.Getenv("KEYCLOAK_ISSUER_URL"),
		KeycloakJWKSURL:   os.Getenv("KEYCLOAK_JWKS_URL"),
		KeycloakAudience:  os.Getenv("KEYCLOAK_AUDIENCE"),

		KeycloakURL:          os.Getenv("KEYCLOAK_URL"),
		KeycloakRealm:        os.Getenv("KEYCLOAK_REALM"),
		KeycloakClientID:     os.Getenv("KEYCLOAK_CLIENT_ID"),
		KeycloakClientSecret: os.Getenv("KEYCLOAK_CLIENT_SECRET"),
	}

	if cfg.DatabaseURL == "" {
		return Config{}, errors.New("DATABASE_URL is required")
	}
	if cfg.KeycloakIssuerURL == "" {
		return Config{}, errors.New("KEYCLOAK_ISSUER_URL is required")
	}
	if cfg.KeycloakJWKSURL == "" {
		return Config{}, errors.New("KEYCLOAK_JWKS_URL is required")
	}
	if cfg.KeycloakAudience == "" {
		return Config{}, errors.New("KEYCLOAK_AUDIENCE is required")
	}
	if cfg.KeycloakURL == "" {
		return Config{}, errors.New("KEYCLOAK_URL is required")
	}
	if cfg.KeycloakRealm == "" {
		return Config{}, errors.New("KEYCLOAK_REALM is required")
	}
	if cfg.KeycloakClientID == "" {
		return Config{}, errors.New("KEYCLOAK_CLIENT_ID is required")
	}
	if cfg.KeycloakClientSecret == "" {
		return Config{}, errors.New("KEYCLOAK_CLIENT_SECRET is required")
	}

	return cfg, nil
}

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
