package api

import (
	"net/http"

	"gorm.io/gorm"

	"parking/internal/auth"
)

func NewRouter(db *gorm.DB, authenticator *auth.Authenticator) http.Handler {
	mux := http.NewServeMux()

	healthHandler := HealthHandler{DB: db}
	mux.HandleFunc("GET /api/health", healthHandler.Check)

	currentUserHandler := auth.RequireAnyRole(http.HandlerFunc(GetCurrentUser), "user")
	mux.Handle("GET /api/me", authenticator.RequireUser(currentUserHandler))

	return mux
}
