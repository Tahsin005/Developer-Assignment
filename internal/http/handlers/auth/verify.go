package auth

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/soyaibzihad10/Developer-Assignment/internal/database"
	"github.com/soyaibzihad10/Developer-Assignment/internal/http/utils"
)

func VerifyEmailHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Token is required", http.StatusBadRequest)
		return
	}

	// Debug logging
	log.Printf("Received verification request with token length: %d", len(token))

	claims, err := utils.ValidateEmailVerificationToken(token)
	if err != nil {
		log.Printf("Invalid verification token: %v", err)
		http.Error(w, "Invalid or expired verification token", http.StatusBadRequest)
		return
	}

	userID := claims.UserID
	email := claims.Email

	// Update user verification status
	if err := database.UpdateEmailVerificationStatus(database.DB, userID); err != nil {
		log.Printf("Error updating verification status: %v", err)
		http.Error(w, "Verification failed", http.StatusBadRequest)
		return
	}

	log.Printf("Successfully verified email for user: %s", email)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Email verified successfully",
	})
}
