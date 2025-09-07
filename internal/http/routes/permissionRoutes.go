package routes

import (
	"github.com/gorilla/mux"
	"github.com/soyaibzihad10/Developer-Assignment/internal/http/handlers/permissions"
	"github.com/soyaibzihad10/Developer-Assignment/internal/http/middleware"
)

func RegisterPermissionRoutes(r *mux.Router) {
	protected := r.PathPrefix("/api/v1/permissions").Subrouter()
	protected.Use(middleware.IsAuthenticated)
	protected.Use(middleware.RequireAdmin)

	protected.HandleFunc("", permissions.ListPermissionsHandler).Methods("GET")
	protected.HandleFunc("/{permission_id}", permissions.GetPermissionHandler).Methods("GET")
}
