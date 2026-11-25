package keycloak

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path"

	"golang.org/x/oauth2/clientcredentials"
)

var ErrUserAlreadyExists = errors.New("keycloak user already exists")

type NewUser struct {
	Email             string
	FirstName         string
	LastName          string
	TemporaryPassword string
	// The user gets the realm roles of this group, e.g. "merchants" gives the "merchant" role.
	Group string
}

// AdminClient calls the Keycloak Admin API as the psp-backend service account.
type AdminClient struct {
	realmAdminURL string
	httpClient    *http.Client
}

func NewAdminClient(keycloakURL string, realm string, clientID string, clientSecret string) *AdminClient {
	credentials := clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     keycloakURL + "/realms/" + realm + "/protocol/openid-connect/token",
	}

	// This HTTP client gets a service account token and renews it when it expires.
	return &AdminClient{
		realmAdminURL: keycloakURL + "/admin/realms/" + realm,
		httpClient:    credentials.Client(context.Background()),
	}
}

// CreateUser creates a user who must change the temporary password and set up TOTP on first login.
// It returns the Keycloak user ID.
func (c *AdminClient) CreateUser(ctx context.Context, user NewUser) (string, error) {
	body := map[string]any{
		"username":      user.Email,
		"email":         user.Email,
		"firstName":     user.FirstName,
		"lastName":      user.LastName,
		"enabled":       true,
		"emailVerified": true,
		"credentials": []map[string]any{
			{"type": "password", "value": user.TemporaryPassword, "temporary": true},
		},
		"requiredActions": []string{"UPDATE_PASSWORD", "CONFIGURE_TOTP"},
		"groups":          []string{"/" + user.Group},
	}

	response, err := c.send(ctx, http.MethodPost, "/users", body)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusConflict {
		return "", ErrUserAlreadyExists
	}
	if response.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("create user: keycloak returned %d", response.StatusCode)
	}

	// Keycloak returns the new user's URL, ending with the user ID.
	return path.Base(response.Header.Get("Location")), nil
}

func (c *AdminClient) DeleteUser(ctx context.Context, userID string) error {
	response, err := c.send(ctx, http.MethodDelete, "/users/"+userID, nil)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("delete user: keycloak returned %d", response.StatusCode)
	}
	return nil
}

func (c *AdminClient) send(ctx context.Context, method string, urlPath string, body any) (*http.Response, error) {
	var requestBody bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&requestBody).Encode(body); err != nil {
			return nil, err
		}
	}

	request, err := http.NewRequestWithContext(ctx, method, c.realmAdminURL+urlPath, &requestBody)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")

	return c.httpClient.Do(request)
}
