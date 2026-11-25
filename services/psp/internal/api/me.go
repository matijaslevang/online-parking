package api

import (
	"encoding/json"
	"net/http"

	"psp/internal/auth"
)

func GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.UserFromContext(r.Context())

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
