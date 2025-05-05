package routes

import (
	"github.com/gorilla/mux"
)

func RegisterRoutes() *mux.Router {
	r := mux.NewRouter()

	api := r.PathPrefix("/api").Subrouter()
	v1 := api.PathPrefix("/v1").Subrouter()

	RegisterCheckHealthRoutes(v1)
	RegisterAuthRoutes(v1)
	return r
}