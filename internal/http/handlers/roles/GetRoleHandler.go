package roles

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/soyaibzihad10/Developer-Assignment/internal/database"
)

func GetRoleHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["role_id"]
	role, err := database.GetRoleByID(id)
	if err != nil {
		http.Error(w, "Role not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(role)
}
