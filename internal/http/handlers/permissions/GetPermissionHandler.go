package permissions

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/soyaibzihad10/Developer-Assignment/internal/database"
)

func GetPermissionHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["permission_id"]

	permission, err := database.GetPermissionByID(id)
	if err != nil {
		http.Error(w, "Permission not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":     "success",
		"permission": permission,
	})
}
