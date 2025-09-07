package auth

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/soyaibzihad10/Developer-Assignment/internal/http/middleware"
)

func CheckAuthHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    // Get values from context
    userID := r.Context().Value(middleware.UserIDKey)
    userType := r.Context().Value(middleware.UserTypeKey)
    userEmail := r.Context().Value(middleware.UserEmailKey)

    // Debug logging
    log.Printf("Raw context values - ID: %v, Type: %v, Email: %v", userID, userType, userEmail)

    if userID == nil || userEmail == nil {
        w.WriteHeader(http.StatusUnauthorized)
        json.NewEncoder(w).Encode(map[string]interface{}{
            "status":        "error",
            "authenticated": false,
            "message":       "User not authenticated",
        })
        return
    }

    // Determine proper user type
    userTypeStr := "user"
    if userType != nil {
        if ut, ok := userType.(string); ok {
            if ut == "system_admin" || ut == "admin" {
                userTypeStr = ut
            }
        }
        log.Printf("User type determined: %s", userTypeStr)
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "status":        "success",
        "authenticated": true,
        "user": map[string]interface{}{
            "id":        userID,
            "email":     userEmail,
            "user_type": userTypeStr,
            "is_admin":  userTypeStr == "system_admin" || userTypeStr == "admin",
        },
    })
}