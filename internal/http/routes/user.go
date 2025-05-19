package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tahsin005/affpilot-auth/internal/database"
	userHandlers "github.com/tahsin005/affpilot-auth/internal/http/handlers/user"
	"github.com/tahsin005/affpilot-auth/internal/http/middleware"
)

func RegisterUserRoutes(r *mux.Router, secretKey string) {
	users := r.PathPrefix("/users").Subrouter()
	users.Use(middleware.AuthMiddleware(secretKey))

	users.Handle("/", middleware.PermissionMiddleware(database.DB, "user:read:all")(http.HandlerFunc(userHandlers.UsersListHandler))).Methods(http.MethodGet)

	users.Handle("/{user_id}", middleware.SelfOrAuthorizedMiddleware(database.DB, "user:read:all")(http.HandlerFunc(userHandlers.UserDetailsHandler))).Methods(http.MethodGet)

	users.Handle("/{user_id}", middleware.SelfOrAuthorizedMiddleware(database.DB, "user:update:all")(http.HandlerFunc(userHandlers.UserUpdateDetailsHandler))).Methods(http.MethodPut)

	users.Handle("/{user_id}/request-deletion", middleware.SelfOnlyMiddleware()(http.HandlerFunc(userHandlers.UserDeletionRequestHandler))).Methods(http.MethodPost)

	users.Handle("/{user_id}", middleware.RoleMiddleware("system_admin", "admin", "moderator")(http.HandlerFunc(userHandlers.UserDeleteHandler))).Methods(http.MethodDelete)

	users.Handle("/{user_id}/role", middleware.RoleMiddleware("system_admin", "admin")(http.HandlerFunc(userHandlers.UserRoleChangeHandler))).Methods(http.MethodPost)

	users.Handle("/{user_id}/promote/admin", middleware.PermissionMiddleware(database.DB, "user:promote:admin")(http.HandlerFunc(userHandlers.UserPromoteAdminHandler))).Methods(http.MethodPost)

	users.Handle("/{user_id}/promote/moderator", middleware.PermissionMiddleware(database.DB, "user:promote:moderator")(http.HandlerFunc(userHandlers.UserPromoteModeratorHandler))).Methods(http.MethodPost)

	users.Handle("/{user_id}/demote", middleware.PermissionMiddleware(database.DB, "user:demote")(http.HandlerFunc(userHandlers.UserRoleDemoteHandler))).Methods(http.MethodPost)
	
	r.Handle("/me", middleware.AuthMiddleware(secretKey)(http.HandlerFunc(userHandlers.GetCurrentUser))).Methods(http.MethodGet)
	r.Handle("/me/permissions", middleware.AuthMiddleware(secretKey)(http.HandlerFunc(userHandlers.CurrentUserPermissions))).Methods(http.MethodGet)
}
