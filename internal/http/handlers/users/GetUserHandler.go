package users

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/soyaibzihad10/Developer-Assignment/internal/database"
	"github.com/soyaibzihad10/Developer-Assignment/internal/http/middleware"
)

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	requestorID := r.Context().Value(middleware.UserIDKey).(string)
	requestorType := r.Context().Value(middleware.UserTypeKey).(string)

	if userID != requestorID && requestorType == "user" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	user, err := database.GetUserByID(database.DB, userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"user":   user,
	})
}
