package users

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/soyaibzihad10/Developer-Assignment/internal/database"
)

func PromoteToAdminHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	if err := database.PromoteUserToAdmin(database.DB, userID); err != nil {
		http.Error(w, "Failed to promote user to admin", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "User promoted to admin successfully",
	})
}
