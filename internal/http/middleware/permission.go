package middleware

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
)

func PermissionMiddleware(db *sql.DB, requiredPermission string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value("user").(*Claims)
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			var hasPermission bool
			query := `
				SELECT EXISTS (
					SELECT 1
					FROM user_roles ur
					JOIN role_permissions rp ON ur.role_id = rp.role_id
					JOIN permissions p ON rp.permission_id = p.id
					WHERE ur.user_id = $1
					AND p.name = $2
				)
			`
			err := db.QueryRow(query, claims.UserID, requiredPermission).Scan(&hasPermission)
			if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			if !hasPermission {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func SelfOrAuthorizedMiddleware(db *sql.DB, requiredPermission string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value("user").(*Claims)
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			vars := mux.Vars(r)
			targetUserID := vars["user_id"]

			if targetUserID == claims.UserID {
				next.ServeHTTP(w, r)
				return
			}

			var hasPermission bool
			query := `
				SELECT EXISTS (
					SELECT 1
					FROM user_roles ur
					JOIN role_permissions rp ON ur.role_id = rp.role_id
					JOIN permissions p ON rp.permission_id = p.id
					WHERE ur.user_id = $1
					AND p.name = $2
				)
			`
			err := db.QueryRow(query, claims.UserID, requiredPermission).Scan(&hasPermission)
			if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			if !hasPermission {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}