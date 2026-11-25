package api

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/mail"
	"strings"

	"gorm.io/gorm"

	"psp/internal/auth"
	"psp/internal/keycloak"
	"psp/internal/models"
)

type MerchantHandler struct {
	DB       *gorm.DB
	Keycloak *keycloak.AdminClient
}

type createMerchantRequest struct {
	Name             string `json:"name"`
	Email            string `json:"email"`
	ContactFirstName string `json:"contactFirstName"`
	ContactLastName  string `json:"contactLastName"`
}

type createMerchantResponse struct {
	Merchant          models.Merchant `json:"merchant"`
	TemporaryPassword string          `json:"temporaryPassword"`
}

func (h MerchantHandler) List(w http.ResponseWriter, r *http.Request) {
	var merchants []models.Merchant
	if err := h.DB.Order("created_at DESC").Find(&merchants).Error; err != nil {
		writeError(w, http.StatusInternalServerError, "could not load merchants")
		return
	}
	writeJSON(w, http.StatusOK, merchants)
}

// Create registers a merchant and its login account. The temporary password is returned only once
// and is never stored by the PSP. The merchant must change it on first login.
func (h MerchantHandler) Create(w http.ResponseWriter, r *http.Request) {
	var request createMerchantRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	request.Name = strings.TrimSpace(request.Name)
	request.Email = strings.ToLower(strings.TrimSpace(request.Email))
	request.ContactFirstName = strings.TrimSpace(request.ContactFirstName)
	request.ContactLastName = strings.TrimSpace(request.ContactLastName)

	if request.Name == "" || request.ContactFirstName == "" || request.ContactLastName == "" {
		writeError(w, http.StatusBadRequest, "name, contact first name and contact last name are required")
		return
	}
	if _, err := mail.ParseAddress(request.Email); err != nil {
		writeError(w, http.StatusBadRequest, "invalid email")
		return
	}

	var existingCount int64
	err := h.DB.Model(&models.Merchant{}).Where("name = ? OR email = ?", request.Name, request.Email).Count(&existingCount).Error
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not check existing merchants")
		return
	}
	if existingCount > 0 {
		writeError(w, http.StatusConflict, "merchant with this name or email already exists")
		return
	}

	temporaryPassword := generateTemporaryPassword()
	keycloakUserID, err := h.Keycloak.CreateUser(r.Context(), keycloak.NewUser{
		Email:             request.Email,
		FirstName:         request.ContactFirstName,
		LastName:          request.ContactLastName,
		TemporaryPassword: temporaryPassword,
		Group:             "merchants",
	})
	if errors.Is(err, keycloak.ErrUserAlreadyExists) {
		writeError(w, http.StatusConflict, "login account with this email already exists")
		return
	}
	if err != nil {
		log.Printf("create merchant: %v", err)
		writeError(w, http.StatusInternalServerError, "could not create login account")
		return
	}

	superAdmin, _ := auth.UserFromContext(r.Context())
	merchant := models.Merchant{
		Name:             request.Name,
		Email:            request.Email,
		ContactFirstName: request.ContactFirstName,
		ContactLastName:  request.ContactLastName,
		KeycloakUserID:   keycloakUserID,
		CreatedByUserID:  superAdmin.ID,
	}
	if err := h.DB.Create(&merchant).Error; err != nil {
		log.Printf("create merchant: save: %v", err)
		h.deleteKeycloakUser(r, keycloakUserID)
		writeError(w, http.StatusInternalServerError, "could not save merchant")
		return
	}

	log.Printf("merchant %d (%s) registered by user %s", merchant.ID, merchant.Name, superAdmin.ID)
	writeJSON(w, http.StatusCreated, createMerchantResponse{Merchant: merchant, TemporaryPassword: temporaryPassword})
}

// Keeps Keycloak and the PSP database consistent when registration fails halfway.
func (h MerchantHandler) deleteKeycloakUser(r *http.Request, keycloakUserID string) {
	if err := h.Keycloak.DeleteUser(r.Context(), keycloakUserID); err != nil {
		log.Printf("create merchant: could not delete keycloak user %s: %v", keycloakUserID, err)
	}
}

// rand.Text returns 26 random characters (A-Z, 2-7). Lowercased, it must still contain
// a digit and a letter to satisfy the realm password policy, so retry in the rare case it does not.
func generateTemporaryPassword() string {
	for {
		password := strings.ToLower(rand.Text())
		hasDigit := strings.ContainsAny(password, "234567")
		hasLetter := strings.ContainsAny(password, "abcdefghijklmnopqrstuvwxyz")
		if hasDigit && hasLetter {
			return password
		}
	}
}
