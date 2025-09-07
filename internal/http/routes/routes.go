package routes

import "github.com/gorilla/mux"

func SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	RegisterPublicRoutes(r)
	RegisterAuthRoutes(r)
	RegisterUserRoutes(r)
	RegisterRoleRoutes(r)
	RegisterPermissionRoutes(r)

	return r
}
