package auth

import (
    "encoding/json"
    "net/http"
    "time"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
    http.SetCookie(w, &http.Cookie{
        Name:     "auth_token",
        Value:    "",
        Path:     "/",
        Expires:  time.Now().Add(-1 * time.Hour), 
        Secure:   true,
    })

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "status":  "success",
        "message": "Successfully logged out",
    })
}