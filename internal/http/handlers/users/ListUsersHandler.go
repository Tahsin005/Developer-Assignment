package users

import (
	"encoding/json"
	"net/http"
	"github.com/soyaibzihad10/Developer-Assignment/internal/database"
)

func ListUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := database.ListUsers()
	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"users":  users,
		"total":  len(users),
	})
}
