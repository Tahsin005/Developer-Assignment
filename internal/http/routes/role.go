package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tahsin005/affpilot-auth/internal/database"
	roleHandlers "github.com/tahsin005/affpilot-auth/internal/http/handlers/role"
	"github.com/tahsin005/affpilot-auth/internal/http/middleware"
)

func RegisterRoleRoutes(r *mux.Router, secretKey string) {
	roles := r.PathPrefix("/roles").Subrouter()
	roles.Use(middleware.AuthMiddleware(secretKey))

	roles.Handle("/", middleware.PermissionMiddleware(database.DB, "role:read")(http.HandlerFunc(roleHandlers.RoleListHandler))).Methods(http.MethodGet)
	roles.Handle("/{role_id}", middleware.PermissionMiddleware(database.DB, "role:read")(http.HandlerFunc(roleHandlers.RoleDetailsHandler))).Methods(http.MethodGet)
	roles.Handle("/", middleware.PermissionMiddleware(database.DB, "role:create")(http.HandlerFunc(roleHandlers.RoleCreateHandler))).Methods(http.MethodPost)
	roles.Handle("/{role_id}", middleware.PermissionMiddleware(database.DB, "role:update")(http.HandlerFunc(roleHandlers.RoleUpdateHandler))).Methods(http.MethodPut)
	roles.Handle("/{role_id}", middleware.PermissionMiddleware(database.DB, "role:delete")(http.HandlerFunc(roleHandlers.RoleDeleteHandler))).Methods(http.MethodDelete)
}