package routes

import (
	"github.com/gorilla/mux"
	"github.com/soyaibzihad10/Developer-Assignment/internal/http/handlers/roles"
	"github.com/soyaibzihad10/Developer-Assignment/internal/http/middleware"
)

func RegisterRoleRoutes(r *mux.Router) {
	protected := r.PathPrefix("/api/v1/roles").Subrouter()
	protected.Use(middleware.IsAuthenticated)
	protected.Use(middleware.RequireAdmin)

	protected.HandleFunc("", roles.ListRolesHandler).Methods("GET")
	protected.HandleFunc("/{role_id}", roles.GetRoleHandler).Methods("GET")
	protected.HandleFunc("", roles.CreateRoleHandler).Methods("POST")
	protected.HandleFunc("/{role_id}", roles.UpdateRoleHandler).Methods("PUT")
	protected.HandleFunc("/{role_id}", roles.DeleteRoleHandler).Methods("DELETE")
}
