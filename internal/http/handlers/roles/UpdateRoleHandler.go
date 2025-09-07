package roles

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/soyaibzihad10/Developer-Assignment/internal/database"
	"github.com/soyaibzihad10/Developer-Assignment/internal/models"
)

func UpdateRoleHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["role_id"]
	var role models.Role
	if err := json.NewDecoder(r.Body).Decode(&role); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := database.UpdateRole(id, &role); err != nil {
		if err == database.ErrNotFound {
			http.Error(w, "Role not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to update role", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(role)
}
