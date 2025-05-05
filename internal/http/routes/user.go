package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tahsin005/affpilot-auth/internal/database"
	"github.com/tahsin005/affpilot-auth/internal/http/handlers"
	"github.com/tahsin005/affpilot-auth/internal/http/middleware"
)

func RegisterUserRoutes(r *mux.Router, secretKey string) {
	users := r.PathPrefix("/users").Subrouter()
	users.Use(middleware.AuthMiddleware(secretKey))

	users.Handle("/", middleware.PermissionMiddleware(database.DB, "user:read:all")(http.HandlerFunc(handlers.CheckHealth))).Methods(http.MethodGet)
	users.Handle("/{user_id}", middleware.SelfOrAuthorizedMiddleware(database.DB, "user:read:all")(http.HandlerFunc(handlers.CheckHealth))).Methods(http.MethodGet)
	users.Handle("/{user_id}", middleware.SelfOrAuthorizedMiddleware(database.DB, "user:update:all")(http.HandlerFunc(handlers.CheckHealth))).Methods(http.MethodPut)
	users.Handle("/{user_id}/request-deletion", middleware.SelfOrAuthorizedMiddleware(database.DB, "user:delete:self")(http.HandlerFunc(handlers.CheckHealth))).Methods(http.MethodPost)
	users.Handle("/{user_id}", middleware.PermissionMiddleware(database.DB, "user:delete:all")(http.HandlerFunc(handlers.CheckHealth))).Methods(http.MethodDelete)
	users.Handle("/{user_id}/role", middleware.PermissionMiddleware(database.DB, "user:promote:admin")(http.HandlerFunc(handlers.CheckHealth))).Methods(http.MethodPost)
	users.Handle("/{user_id}/promote/admin", middleware.PermissionMiddleware(database.DB, "user:promote:admin")(http.HandlerFunc(handlers.CheckHealth))).Methods(http.MethodPost)
	users.Handle("/{user_id}/promote/moderator", middleware.PermissionMiddleware(database.DB, "user:promote:moderator")(http.HandlerFunc(handlers.CheckHealth))).Methods(http.MethodPost)
	users.Handle("/{user_id}/demote", middleware.PermissionMiddleware(database.DB, "user:demote")(http.HandlerFunc(handlers.CheckHealth))).Methods(http.MethodPost)
	
	r.Handle("/me", middleware.AuthMiddleware(secretKey)(http.HandlerFunc(handlers.CheckHealth))).Methods(http.MethodGet)
	r.Handle("/me/permissions", middleware.AuthMiddleware(secretKey)(http.HandlerFunc(handlers.CheckHealth))).Methods(http.MethodGet)
}
