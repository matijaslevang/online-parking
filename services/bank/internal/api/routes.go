package api

import (
	"net/http"

	"gorm.io/gorm"
)

func NewRouter(db *gorm.DB) http.Handler {
	mux := http.NewServeMux()

	healthHandler := HealthHandler{DB: db}
	mux.HandleFunc("GET /api/health", healthHandler.Check)

	return mux
}
