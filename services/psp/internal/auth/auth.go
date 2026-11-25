package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"slices"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
)

type User struct {
	ID       string   `json:"id"`
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Name     string   `json:"name"`
	Roles    []string `json:"roles"`
}

type tokenClaims struct {
	Username    string `json:"preferred_username"`
	Email       string `json:"email"`
	Name        string `json:"name"`
	RealmAccess struct {
		Roles []string `json:"roles"`
	} `json:"realm_access"`
}

type contextKey struct{}

type Authenticator struct {
	verifier *oidc.IDTokenVerifier
}

// The verifier checks the signature (keys downloaded from jwksURL), issuer, audience and expiry.
// Despite its name, IDTokenVerifier works for any Keycloak JWT, including access tokens.
func NewAuthenticator(issuerURL string, jwksURL string, audience string) *Authenticator {
	keySet := oidc.NewRemoteKeySet(context.Background(), jwksURL)
	verifier := oidc.NewVerifier(issuerURL, keySet, &oidc.Config{ClientID: audience})
	return &Authenticator{verifier: verifier}
}

func (a *Authenticator) RequireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawToken, found := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
		if !found {
			writeError(w, http.StatusUnauthorized, "missing bearer token")
			return
		}

		token, err := a.verifier.Verify(r.Context(), rawToken)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "invalid token")
			return
		}

		var claims tokenClaims
		if err := token.Claims(&claims); err != nil {
			writeError(w, http.StatusUnauthorized, "invalid token claims")
			return
		}

		user := User{
			ID:       token.Subject,
			Username: claims.Username,
			Email:    claims.Email,
			Name:     claims.Name,
			Roles:    claims.RealmAccess.Roles,
		}
		ctx := context.WithValue(r.Context(), contextKey{}, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireAnyRole must be used inside RequireUser.
func RequireAnyRole(next http.Handler, roles ...string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, found := UserFromContext(r.Context())
		if !found {
			writeError(w, http.StatusUnauthorized, "not authenticated")
			return
		}

		for _, role := range roles {
			if slices.Contains(user.Roles, role) {
				next.ServeHTTP(w, r)
				return
			}
		}
		writeError(w, http.StatusForbidden, "missing required role")
	})
}

func UserFromContext(ctx context.Context) (User, bool) {
	user, found := ctx.Value(contextKey{}).(User)
	return user, found
}

func writeError(w http.ResponseWriter, statusCode int, message string) {
	if statusCode == http.StatusUnauthorized {
		w.Header().Set("WWW-Authenticate", "Bearer")
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
