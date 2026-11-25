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
}

func Load() (Config, error) {
	cfg := Config{
		Port:              getEnv("PORT", "8080"),
		DatabaseURL:       os.Getenv("DATABASE_URL"),
		KeycloakIssuerURL: os.Getenv("KEYCLOAK_ISSUER_URL"),
		KeycloakJWKSURL:   os.Getenv("KEYCLOAK_JWKS_URL"),
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

	return cfg, nil
}

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
