package middleware

import (
	"net/http"

	"github.com/gorilla/mux"
)

func RequireAdminModeratorSelf(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userType := r.Context().Value(UserTypeKey).(string)
		userID := r.Context().Value(UserIDKey).(string)

		// Check if request is for a specific user
		vars := mux.Vars(r)
		requestedUserID := vars["user_id"]

		// Allow if admin/system_admin OR if user accessing their own data
		isAdmin := userType == "admin" || userType == "system_admin" || userType == "moderator"
		isSelfAccess := requestedUserID != "" && requestedUserID == userID

		if !isAdmin && !isSelfAccess {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
