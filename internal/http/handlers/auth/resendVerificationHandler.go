package auth

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/soyaibzihad10/Developer-Assignment/internal/database"
	"github.com/soyaibzihad10/Developer-Assignment/internal/email"
	"github.com/soyaibzihad10/Developer-Assignment/internal/http/utils"
)

type ResendVerificationRequest struct {
	Email string `json:"email"`
}

// ResendVerificationHandler handles resending of verification emails
func ResendVerificationHandler(w http.ResponseWriter, r *http.Request) {
	// Set JSON content type early
	w.Header().Set("Content-Type", "application/json")

	var req ResendVerificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSONError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Email == "" {
		sendJSONError(w, "Email is required", http.StatusBadRequest)
		return
	}

	// Fetch user from database
	user, err := database.GetUserByEmail(database.DB, req.Email)
	if err != nil {
		sendJSONError(w, "User not found", http.StatusNotFound)
		return
	}

	// Check if email is already verified
	if user.EmailVerified {
		sendJSONError(w, "Email is already verified", http.StatusBadRequest)
		return
	}

	// Generate new verification token
	token, err := utils.GenerateEmailVerificationToken(user.ID, user.Email)
	if err != nil {
		log.Printf("Error generating verification token: %v", err)
		sendJSONError(w, "Could not generate verification token", http.StatusInternalServerError)
		return
	}

	// Update user with new token and expiry
	user.VerificationToken = token
	user.TokenExpiry = time.Now().Add(5 * time.Minute)

	if err := database.UpdateUserVerificationToken(database.DB, user); err != nil {
		log.Printf("Error updating verification token: %v", err)
		sendJSONError(w, "Could not update verification token", http.StatusInternalServerError)
		return
	}

	// Send verification email in background
	go func() {
		if err := email.SendVerificationEmail(user.Email, user.VerificationToken); err != nil {
			log.Printf("Failed to send verification email to %s: %v", user.Email, err)
		}
	}()

	// Return success response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Verification email resent successfully",
	})
}

// Helper function for JSON error responses
func sendJSONError(w http.ResponseWriter, message string, status int) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "error",
		"message": message,
	})
}
