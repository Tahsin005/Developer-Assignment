package auth

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"path"

	"github.com/soyaibzihad10/Developer-Assignment/internal/config"
	"github.com/soyaibzihad10/Developer-Assignment/internal/database"
	"github.com/soyaibzihad10/Developer-Assignment/internal/email"
	"github.com/soyaibzihad10/Developer-Assignment/internal/http/utils"
)

type PasswordResetRequest struct {
	Email string `json:"email"`
}

// RequestPasswordResetHandler handles password reset requests
func RequestPasswordResetHandler(w http.ResponseWriter, r *http.Request) {
	cfg := config.GetConfig()
	var req PasswordResetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Always set JSON content type
	w.Header().Set("Content-Type", "application/json")

	user, err := database.GetUserByEmail(database.DB, req.Email)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "success",
			"message": "If your email is registered, you will receive a password reset link",
		})
		return
	}

	// Generate JWT reset token
	resetToken, err := utils.GeneratePasswordResetToken(user.ID, user.Email)
	if err != nil {
		log.Printf("Error generating reset token: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Construct reset URL
	resetURL, err := url.Parse(cfg.App.AppURL)
	if err != nil {
		log.Printf("Invalid base URL: %v", err)
		http.Error(w, "Server configuration error", http.StatusInternalServerError)
		return
	}

	resetURL.Path = path.Join("auth", "reset-password")
	q := resetURL.Query()
	q.Set("token", resetToken)
	resetURL.RawQuery = q.Encode()

	// Send email in goroutine
	go func() {
		if err := email.SendPasswordResetEmail(user.Email, resetURL.String(), resetToken); err != nil {
			log.Printf("Failed to send reset email to %s: %v", user.Email, err)
		}
	}()

	// Return success response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "If your email is registered, you will receive a password reset link",
	})
}
