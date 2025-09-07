package permissions

import (
	"encoding/json"
	"net/http"

	"github.com/soyaibzihad10/Developer-Assignment/internal/database"
)

func ListPermissionsHandler(w http.ResponseWriter, r *http.Request) {
	permissions, err := database.ListPermissions()
	if err != nil {
		http.Error(w, "Failed to fetch permissions", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":      "success",
		"permissions": permissions,
		"total":       len(permissions),
	})
}
