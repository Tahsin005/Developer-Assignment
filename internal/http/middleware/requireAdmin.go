package middleware

import (
	"log"
	"net/http"
)

func RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get user type from context
		userType, ok := r.Context().Value(UserTypeKey).(string)
		if !ok {
			log.Printf("Admin check failed: No user type in context")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Check if user is admin or system_admin
		if userType != "admin" && userType != "system_admin" {
			log.Printf("Admin check failed: User type '%s' is not admin", userType)
			http.Error(w, "Forbidden - Admin access required", http.StatusForbidden)
			return
		}

		log.Printf("Admin access granted for user type: %s", userType)
		next.ServeHTTP(w, r)
	})
}
