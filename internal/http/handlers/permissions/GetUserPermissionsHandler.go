package permissions

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/soyaibzihad10/Developer-Assignment/internal/database"
	"github.com/soyaibzihad10/Developer-Assignment/internal/http/middleware"
)

func GetUserPermissionsHandler(w http.ResponseWriter, r *http.Request) {
	userType := r.Context().Value(middleware.UserTypeKey).(string)
	if userType == "" {
		http.Error(w, "User type not found", http.StatusUnauthorized)
		return
	}

	userID := r.Context().Value(middleware.UserIDKey).(string)
	permissions, err := database.GetUserPermissions(userID)
	if err != nil {
		http.Error(w, "Failed to fetch user permissions", http.StatusInternalServerError)
		return
	}

	fmt.Printf("chill %s", userID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":      "success",
		"permissions": permissions,
		"total":       len(permissions),
	})
}
