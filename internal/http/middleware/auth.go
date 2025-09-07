package middleware

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v4"
)

type contextKey string

const (
    UserIDKey    contextKey = "user_id"
    UserTypeKey  contextKey = "user_type"
    UserEmailKey contextKey = "user_email"
)

// IsAuthenticated checks if the user is authenticated by verifying the JWT token in the cookie
func IsAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get token from cookie
		cookie, err := r.Cookie("auth_token")
		if err != nil {
			log.Printf("Auth failed: No auth_token cookie found - %v", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Parse token with claims
		token, err := jwt.ParseWithClaims(cookie.Value, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			log.Printf("Auth failed: Invalid token - %v", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Extract and validate claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			log.Printf("Auth failed: Invalid claims format")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Extract user information with type assertions
		userID, _ := claims["user_id"].(string)
		userType, _ := claims["user_type"].(string)
		userEmail, _ := claims["email"].(string)

		if userID == "" || userEmail == "" {
			log.Printf("Auth failed: Missing required claims")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Set default user type if not present
		if userType == "" {
			userType = "user"
		}

		// Create context with user information
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		ctx = context.WithValue(ctx, UserTypeKey, userType)
		ctx = context.WithValue(ctx, UserEmailKey, userEmail)

		log.Printf("Auth success - UserID: %s, Type: %s, Email: %s", userID, userType, userEmail)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
