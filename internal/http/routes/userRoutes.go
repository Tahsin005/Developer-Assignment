package routes

import (
	"github.com/gorilla/mux"
	"github.com/soyaibzihad10/Developer-Assignment/internal/http/handlers/users"
	"github.com/soyaibzihad10/Developer-Assignment/internal/http/handlers/permissions"
	"github.com/soyaibzihad10/Developer-Assignment/internal/http/middleware"
)

func RegisterUserRoutes(r *mux.Router) {
	protected := r.PathPrefix("/api/v1").Subrouter()
	protected.Use(middleware.IsAuthenticated)

	protected.HandleFunc("/me", users.GetProfileHandler).Methods("GET")
	protected.HandleFunc("/me/permissions", permissions.GetUserPermissionsHandler).Methods("GET")

	userRoutes := protected.PathPrefix("/users").Subrouter()

	adminList := userRoutes.PathPrefix("").Subrouter()
	adminList.Use(middleware.RequireAdmin)
	adminList.HandleFunc("", users.ListUsersHandler).Methods("GET")
	adminList.HandleFunc("/{user_id}", users.DeleteUserHandler).Methods("DELETE")

	adminSelf := userRoutes.PathPrefix("").Subrouter()
	adminSelf.Use(middleware.RequireAdminModeratorSelf)
	adminSelf.HandleFunc("/{user_id}", users.GetUserHandler).Methods("GET")

	self := userRoutes.PathPrefix("").Subrouter()
	self.Use(middleware.RequireAdminOrSelf)
	self.HandleFunc("/{user_id}", users.UpdateUserHandler).Methods("PUT")
	self.HandleFunc("/{user_id}/request-deletion", users.RequestDeletionHandler).Methods("POST")

	roleChange := userRoutes.PathPrefix("").Subrouter()
	roleChange.Use(middleware.RequireAdmin)
	roleChange.HandleFunc("/{user_id}/role", users.ChangeRoleHandler).Methods("POST")

	admin := userRoutes.PathPrefix("").Subrouter()
	admin.Use(middleware.RequireAdmin)
	admin.HandleFunc("/{user_id}/promote/moderator", users.PromoteToModeratorHandler).Methods("POST")
	admin.HandleFunc("/{user_id}/demote", users.DemoteUserHandler).Methods("POST")

	sysAdmin := userRoutes.PathPrefix("").Subrouter()
	sysAdmin.Use(middleware.RequireSystemAdmin)
	sysAdmin.HandleFunc("/{user_id}/promote/admin", users.PromoteToAdminHandler).Methods("POST")
}
