package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)


func UserLogoutHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    http.SetCookie(w, &http.Cookie{
        Name:     "token",
        Value:    "",
        Expires:  time.Unix(0, 0),
        HttpOnly: true,
        Secure:   false,
        SameSite: http.SameSiteStrictMode,
        Path:     "/",
    })
    http.SetCookie(w, &http.Cookie{
        Name:     "role",
        Value:    "",
        Expires:  time.Unix(0, 0),
        HttpOnly: true,
        Secure:   false,
        SameSite: http.SameSiteStrictMode,
        Path:     "/",
    })
    http.SetCookie(w, &http.Cookie{
        Name:     "username",
        Value:    "",
        Expires:  time.Unix(0, 0),
        HttpOnly: true,
        Secure:   false,
        SameSite: http.SameSiteStrictMode,
        Path:     "/",
    })
    http.SetCookie(w, &http.Cookie{
        Name:     "id",
        Value:    "",
        Expires:  time.Unix(0, 0),
        HttpOnly: true,
        Secure:   false,
        SameSite: http.SameSiteStrictMode,
        Path:     "/",
    })
    w.WriteHeader(http.StatusOK)
    fmt.Println("Logout succesful")
    json.NewEncoder(w).Encode(map[string]string{
        "message": "Logout successful",
    })
}