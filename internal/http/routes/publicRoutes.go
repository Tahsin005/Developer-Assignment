package routes

import (
	"github.com/gorilla/mux"
	"github.com/soyaibzihad10/Developer-Assignment/internal/http/handlers/health"
)

func RegisterPublicRoutes(r *mux.Router) {
	r.HandleFunc("/ping", health.PingHandler).Methods("GET")
}
