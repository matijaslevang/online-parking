package api

import (
	"net/http"

	"gorm.io/gorm"

	"psp/internal/auth"
	"psp/internal/keycloak"
)

func NewRouter(db *gorm.DB, authenticator *auth.Authenticator, keycloakAdmin *keycloak.AdminClient) http.Handler {
	mux := http.NewServeMux()

	healthHandler := HealthHandler{DB: db}
	mux.HandleFunc("GET /api/health", healthHandler.Check)

	currentUserHandler := auth.RequireAnyRole(http.HandlerFunc(GetCurrentUser), "super-admin", "merchant")
	mux.Handle("GET /api/me", authenticator.RequireUser(currentUserHandler))

	merchantHandler := MerchantHandler{DB: db, Keycloak: keycloakAdmin}
	listMerchantsHandler := auth.RequireAnyRole(http.HandlerFunc(merchantHandler.List), "super-admin")
	createMerchantHandler := auth.RequireAnyRole(http.HandlerFunc(merchantHandler.Create), "super-admin")
	mux.Handle("GET /api/merchants", authenticator.RequireUser(listMerchantsHandler))
	mux.Handle("POST /api/merchants", authenticator.RequireUser(createMerchantHandler))

	return mux
}
