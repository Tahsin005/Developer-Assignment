package middleware

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func RequireAdminOrSelf(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userType := r.Context().Value(UserTypeKey)
		userID := r.Context().Value(UserIDKey)

		log.Println("Middleware Check -- userType:", userType, "userID:", userID)

		vars := mux.Vars(r)
		requestedUserID := vars["user_id"]
		log.Println("Requested user ID:", requestedUserID)

		isAdmin := userType == "admin" || userType == "system_admin"
		isSelfAccess := requestedUserID != "" && requestedUserID == userID

		if !isAdmin && !isSelfAccess {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
