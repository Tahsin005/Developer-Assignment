package roles

import (
	"encoding/json"
	"net/http"

	"github.com/soyaibzihad10/Developer-Assignment/internal/database"
	"github.com/soyaibzihad10/Developer-Assignment/internal/models"
)

func ListRolesHandler(w http.ResponseWriter, r *http.Request) {
	roles, err := database.ListRoles()
	if err != nil {
		http.Error(w, "Failed to fetch roles", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.RoleList{
		Roles: roles,
		Total: len(roles),
	})
}
