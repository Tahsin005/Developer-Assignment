package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/tahsin005/affpilot-auth/internal/database"
	"github.com/tahsin005/affpilot-auth/internal/models"
	"github.com/tahsin005/affpilot-auth/internal/services"
	"github.com/tahsin005/affpilot-auth/internal/utils"
)

func PasswordResetRequestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req models.PasswordResetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid input")
		return
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	if req.Email == "" {
		utils.WriteError(w, http.StatusBadRequest, "Email is required")
		return
	}

	log.Printf("Password reset requested for email: %s", req.Email)

	var userID string
	query := `SELECT id FROM users WHERE email = $1`
	err := database.DB.QueryRow(query, req.Email).Scan(&userID)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("Database error checking email: %v", err)
		utils.WriteError(w, http.StatusInternalServerError, "Database error")
		return
	}

	resetToken := uuid.New().String()
	tokenExpiry := time.Now().UTC().Add(5 * time.Minute)

	if err != sql.ErrNoRows {
		query = `
			UPDATE users
			SET verification_token = $1, token_expiry = $2, updated_at = NOW()
			WHERE id = $3
		`
		_, err = database.DB.Exec(query, resetToken, tokenExpiry, userID)
		if err != nil {
			log.Printf("Failed to update reset token: %v", err)
			utils.WriteError(w, http.StatusInternalServerError, "Failed to process reset request")
			return
		}

		emailBody := fmt.Sprintf(`Hello,
We received a request to reset your Affpilot account password.
Your password reset token is:

		%s

Please use this token to reset your password. The token will expire in 2 hours.
If you did not request this, please ignore this email.
Best regards,  
The Affpilot Team`, resetToken)

		go services.SendEmail(req.Email, emailBody, "Password Reset Request")
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "If an account with that email exists, a password reset token has been sent.",
	})
}
