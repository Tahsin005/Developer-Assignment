package middleware

import (
	"net/http"
)

func RequireSystemAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userType := r.Context().Value(UserTypeKey).(string)

		if userType != "system_admin" {
			http.Error(w, "Forbidden - System Admin access required", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
