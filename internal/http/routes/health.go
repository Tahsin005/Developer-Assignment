package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	healthHandlers "github.com/tahsin005/affpilot-auth/internal/http/handlers/health"
)

func RegisterCheckHealthRoutes(r *mux.Router) {
	r.HandleFunc("/health", healthHandlers.CheckHealth).Methods(http.MethodGet)
}