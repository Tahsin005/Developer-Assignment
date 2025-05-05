package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tahsin005/affpilot-auth/internal/http/handlers"
)

func RegisterAuthRoutes(r *mux.Router) {
	auth := r.PathPrefix("/auth").Subrouter()

	auth.HandleFunc("/register", handlers.UserRegisterHandler).Methods(http.MethodPost)
}