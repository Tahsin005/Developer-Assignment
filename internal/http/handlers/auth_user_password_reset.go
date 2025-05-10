package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/tahsin005/affpilot-auth/internal/database"
	"github.com/tahsin005/affpilot-auth/internal/models"
	"github.com/tahsin005/affpilot-auth/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

func PasswordResetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req models.PasswordResetConfirmRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid input")
		return
	}

	req.Token = strings.TrimSpace(req.Token)
	if req.Token == "" || req.NewPassword == "" {
		utils.WriteError(w, http.StatusBadRequest, "Token and new password are required")
		return
	}

	log.Printf("Password reset attempt with token: %s", req.Token)

	var userID string
	var tokenExpiry time.Time
	query := `
		SELECT id, token_expiry
		FROM users
		WHERE verification_token = $1 AND email_verified = TRUE AND active = TRUE
	`
	err := database.DB.QueryRow(query, req.Token).Scan(&userID, &tokenExpiry)
	if err == sql.ErrNoRows {
		utils.WriteError(w, http.StatusBadRequest, "Invalid or expired reset token")
		return
	}
	if err != nil {
		log.Printf("Database error verifying token: %v", err)
		utils.WriteError(w, http.StatusInternalServerError, "Database error")
		return
	}

	if time.Now().UTC().After(tokenExpiry) {
		utils.WriteError(w, http.StatusBadRequest, "Reset token has expired")
		return
	}
	
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		utils.WriteError(w, http.StatusInternalServerError, "Error processing password")
		return
	}

	query = `
		UPDATE users
		SET password_hash = $1, verification_token = 'reset', token_expiry = NOW(), updated_at = NOW()
		WHERE id = $2
	`
	_, err = database.DB.Exec(query, string(hashedPassword), userID)
	if err != nil {
		log.Printf("Failed to update password: %v", err)
		utils.WriteError(w, http.StatusInternalServerError, "Failed to reset password")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Password reset successfully",
	})
}