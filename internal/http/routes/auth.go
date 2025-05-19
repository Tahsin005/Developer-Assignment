package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	authHandlers "github.com/tahsin005/affpilot-auth/internal/http/handlers/auth"
	"github.com/tahsin005/affpilot-auth/internal/http/middleware"
)

func RegisterAuthRoutes(r *mux.Router, secretKey string) {
	auth := r.PathPrefix("/auth").Subrouter()

	auth.HandleFunc("/register", authHandlers.UserRegisterHandler).Methods(http.MethodPost)
	auth.HandleFunc("/login", authHandlers.UserLoginHandler).Methods(http.MethodPost)
	auth.HandleFunc("/verify/{verification_token}", authHandlers.VerifyEmailHandler).Methods(http.MethodGet)
	auth.HandleFunc("/resend-verification", authHandlers.ResendVerificationEmail).Methods(http.MethodPost)
	auth.HandleFunc("/reset-request", authHandlers.PasswordResetRequestHandler).Methods(http.MethodPost)
	auth.HandleFunc("/password-reset", authHandlers.PasswordResetHandler).Methods(http.MethodPost)

	auth.Handle("/logout", middleware.AuthMiddleware(secretKey)(http.HandlerFunc(authHandlers.UserLogoutHandler))).Methods(http.MethodGet)
}
