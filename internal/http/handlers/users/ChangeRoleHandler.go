package users

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/soyaibzihad10/Developer-Assignment/internal/database"
)

func ChangeRoleHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	var req struct {
		Role string `json:"role"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Role == "admin" || req.Role == "moderator" || req.Role == "system_admin" {
		http.Error(w, "Invalid role, you can't change role to admin or system admin or moderator", http.StatusBadRequest)
		return
	}

	if err := database.ChangeUserRole(database.DB, userID, req.Role); err != nil {
		http.Error(w, "Failed to change role", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Role updated successfully",
	})
}
