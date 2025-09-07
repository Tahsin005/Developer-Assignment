package roles

import (
	"encoding/json"
	"net/http"

	"github.com/soyaibzihad10/Developer-Assignment/internal/database"
	"github.com/soyaibzihad10/Developer-Assignment/internal/models"
)

func CreateRoleHandler(w http.ResponseWriter, r *http.Request) {
	var role models.Role
	if err := json.NewDecoder(r.Body).Decode(&role); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if role.Name == "" {
		http.Error(w, "Role name is required", http.StatusBadRequest)
		return
	}

	if err := database.CreateRole(&role); err != nil {
		http.Error(w, "Failed to create role", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(role)
}
