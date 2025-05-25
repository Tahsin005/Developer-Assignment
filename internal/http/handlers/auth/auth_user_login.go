package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/tahsin005/affpilot-auth/internal/config"
	"github.com/tahsin005/affpilot-auth/internal/database"
	"github.com/tahsin005/affpilot-auth/internal/http/middleware"
	"github.com/tahsin005/affpilot-auth/internal/models"
	"github.com/tahsin005/affpilot-auth/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

func UserLoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid input")
		return
	}

	req.Password = strings.TrimSpace(req.Password)
	req.Username = strings.TrimSpace(req.Username)
	if req.Username == "" || req.Password == "" {
		utils.WriteError(w, http.StatusBadRequest, "Username and password are required")
		return
	}

	var userID uuid.UUID
	var passwordHash, roleName string
	var emailVerified, isActive bool

	query := `
		SELECT u.id, u.password_hash, r.name, u.email_verified, u.active
		FROM users u
		LEFT JOIN user_roles ur ON u.id = ur.user_id
		LEFT JOIN roles r ON ur.role_id = r.id
		WHERE u.username = $1
	`
	err := database.DB.QueryRow(query, req.Username).Scan(&userID, &passwordHash, &roleName, &emailVerified, &isActive)
	if err != nil {
		log.Println(err)
		utils.WriteError(w, http.StatusUnauthorized, "Invalid username or password")
		return
	}

	if !isActive {
		utils.WriteError(w, http.StatusForbidden, "Account is deactivated. Please contact support.")
		return
	}

	if !emailVerified {
		utils.WriteError(w, http.StatusUnauthorized, "Email not verified. Please verify your email before logging in.")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)); err != nil {
		utils.WriteError(w, http.StatusUnauthorized, "Invalid username or password")
		return
	}

	cfg := config.GetConfig()

	claims := &middleware.Claims{
		UserID:   userID.String(),
		Username: req.Username,
		Role:     roleName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(cfg.JWT_SECRET))
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	cookies := []http.Cookie{
		{Name: "token", Value: tokenString, Expires: time.Now().Add(24 * time.Hour), HttpOnly: true, Secure: false, SameSite: http.SameSiteStrictMode, Path: "/"},
		{Name: "role", Value: roleName, Expires: time.Now().Add(1000 * time.Hour), HttpOnly: true, Secure: false, SameSite: http.SameSiteStrictMode, Path: "/"},
		{Name: "username", Value: req.Username, Expires: time.Now().Add(1000 * time.Hour), HttpOnly: true, Secure: false, SameSite: http.SameSiteStrictMode, Path: "/"},
		{Name: "id", Value: userID.String(), Expires: time.Now().Add(1000 * time.Hour), HttpOnly: true, Secure: false, SameSite: http.SameSiteStrictMode, Path: "/"},
	}
	for _, c := range cookies {
		http.SetCookie(w, &c)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":  "Login successful",
		"token":    tokenString,
		"role":     roleName,
		"username": req.Username,
		"id":       userID.String(),
	})
}
