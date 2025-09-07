package users

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/soyaibzihad10/Developer-Assignment/internal/database"
)

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	if err := database.DeleteUser(database.DB, userID); err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
