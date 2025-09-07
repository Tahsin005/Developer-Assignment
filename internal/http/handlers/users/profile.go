package users

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/soyaibzihad10/Developer-Assignment/internal/database"
	"github.com/soyaibzihad10/Developer-Assignment/internal/http/middleware"
)

func GetProfileHandler(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		log.Printf("User ID not found in context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get user profile from database
	user, err := database.GetUserByID(database.DB, userID)
	if err != nil {
		log.Printf("Error fetching user profile: %v", err)
		http.Error(w, "Failed to fetch profile", http.StatusInternalServerError)
		return
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")

	// Return profile without sensitive information
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"user": map[string]interface{}{
			"id":             user.ID,
			"username":       user.Username,
			"email":          user.Email,
			"first_name":     user.FirstName,
			"last_name":      user.LastName,
			"user_type":      user.UserType,
			"email_verified": user.EmailVerified,
			"active":         user.Active,
			"created_at":     user.CreatedAt,
		},
	})
}
