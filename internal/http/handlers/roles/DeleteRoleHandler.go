package roles

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/soyaibzihad10/Developer-Assignment/internal/database"
)

func DeleteRoleHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["role_id"]
	if err := database.DeleteRole(id); err != nil {
		if err == database.ErrNotFound {
			http.Error(w, "Role not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to delete role", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
