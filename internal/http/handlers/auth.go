package handlers

import (
	"encoding/json"
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

func UserRegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid input")
		return
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.Username = strings.TrimSpace(req.Username)

	if req.Email == "" || req.Password == "" || req.Username == "" || req.FirstName == "" || req.LastName == "" {
		http.Error(w, "Username, email, password, first name and last name are required", http.StatusBadRequest)
		return
	}

	// Check if email or username already exists
	var exists bool
	queryStatement := `
		SELECT EXISTS (SELECT 1 FROM users WHERE email = $1 OR username = $2)
	`
	err := database.DB.QueryRow(
		queryStatement,
		req.Email, req.Username,
	).Scan(&exists)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Database error")
		return
	}
	
	if exists {
		utils.WriteError(w, http.StatusConflict, "Email or username already exists")
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Error hashing password")
		return
	}

	userID := uuid.New()
	verificationToken := uuid.New().String()
	tokenExpiry := time.Now().Add(24 * time.Hour)
	now := time.Now()

	// Insert user
	queryStatement = `
		INSERT INTO users (
			id, username, email, password_hash, first_name, last_name, email_verified,
			user_type, verification_token, token_expiry, deletion_requested, active, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, FALSE,
			'user', $7, $8, FALSE, TRUE, $9, $10
		)
	`
	_, err = database.DB.Exec(queryStatement, userID, req.Username, req.Email, string(hashedPassword), req.FirstName, req.LastName, verificationToken, tokenExpiry, now, now)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	// Assign role
	var roleID uuid.UUID
	queryStatement = `
		SELECT id FROM roles WHERE name = 'user'
	`
	err = database.DB.QueryRow(queryStatement).Scan(&roleID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to fetch role")
		return
	}

	queryStatement = `
		INSERT INTO user_roles (user_id, role_id, assigned_by, created_at)
		VALUES ($1, $2, $3, $4)
	`
	_, err = database.DB.Exec(queryStatement, userID, roleID, userID, now)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to assign role")
		return
	}

	// response
	userResp := models.RegisterUserResponse{
		ID:        userID,
		Username:  req.Username,
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		UserType:  "user",
		CreatedAt: now,
		UpdatedAt: now,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User registered successfully",
		"user":    userResp,
	})
}


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
	query := `
		SELECT u.id, u.password_hash, r.name
		FROM users u
		LEFT JOIN user_roles ur ON u.id = ur.user_id
		LEFT JOIN roles r ON ur.role_id = r.id
		WHERE u.username = $1 AND u.active = TRUE
	`
	err := database.DB.QueryRow(query, req.Username).Scan(&userID, &passwordHash, &roleName)
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, "Invalid username or password")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)); err != nil {
		utils.WriteError(w, http.StatusUnauthorized, "Invalid username or password")
		return
	}

	cfg := config.LoadConfig()

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

	http.SetCookie(w, &http.Cookie{
		Name:     "jwt_token",
		Value:    tokenString,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Login successful",
		"token":   tokenString,
	})
}



func UserLogoutHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    http.SetCookie(w, &http.Cookie{
        Name:     "jwt_token",
        Value:    "",
        Expires:  time.Now().Add(-time.Hour),
        HttpOnly: true,
        Secure:   false,
        SameSite: http.SameSiteStrictMode,
        Path:     "/",
    })
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{
        "message": "Logout successful",
    })
}