package auth

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/soyaibzihad10/Developer-Assignment/internal/database"
	"github.com/soyaibzihad10/Developer-Assignment/internal/http/utils"
	"golang.org/x/crypto/bcrypt"
)

type ResetPasswordRequest struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

func ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var req ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request: %v", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Validate password
	if len(req.Password) < 8 {
		http.Error(w, "Password must be at least 8 characters", http.StatusBadRequest)
		return
	}

	// Validate JWT token and get userID
	userID, err := utils.ValidatePasswordResetToken(req.Token)
	if err != nil {
		log.Printf("Invalid reset token: %v", err)
		http.Error(w, "Invalid or expired reset token", http.StatusBadRequest)
		return
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Update password
	updateQuery := `
        UPDATE users 
        SET password_hash = $1
        WHERE id = $2`

	result, err := database.DB.Exec(updateQuery, string(hashedPassword), userID)
	if err != nil {
		log.Printf("Error updating password: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Password reset failed - user not found", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Password reset successfully",
	})
}
