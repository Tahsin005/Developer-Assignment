package routes

import (
	"github.com/gorilla/mux"
	"github.com/soyaibzihad10/Developer-Assignment/internal/http/handlers/auth"
	"github.com/soyaibzihad10/Developer-Assignment/internal/http/middleware"
)

func RegisterAuthRoutes(r *mux.Router) {
	authRoutes := r.PathPrefix("/api/v1/auth").Subrouter()

	// Public routes
	authRoutes.HandleFunc("/register", auth.RegisterHandler).Methods("POST")
	authRoutes.HandleFunc("/login", auth.LoginHandler).Methods("POST")
	authRoutes.HandleFunc("/verify", auth.VerifyEmailHandler).Methods("GET")
	authRoutes.HandleFunc("/resend-verification", auth.ResendVerificationHandler).Methods("POST")
	authRoutes.HandleFunc("/password-reset-request", auth.RequestPasswordResetHandler).Methods("POST")
	authRoutes.HandleFunc("/password-reset", auth.ResetPasswordHandler).Methods("POST")

	// Protected routes
	protected := r.PathPrefix("/api/v1/auth").Subrouter()
	protected.Use(middleware.IsAuthenticated)
	protected.HandleFunc("/check", auth.CheckAuthHandler).Methods("GET")
	protected.HandleFunc("/logout", auth.LogoutHandler).Methods("POST")
}
