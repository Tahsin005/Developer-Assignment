package auth

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/soyaibzihad10/Developer-Assignment/internal/database"
	"github.com/soyaibzihad10/Developer-Assignment/internal/email"
	"github.com/soyaibzihad10/Developer-Assignment/internal/http/utils"
	"github.com/soyaibzihad10/Developer-Assignment/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}

// RegisterHandler handles user registration
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" || req.Username == "" {
		http.Error(w, "Email, password, and username are required", http.StatusBadRequest)
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Could not hash password", http.StatusInternalServerError)
		return
	}

	user := &models.User{
		Email:        req.Email,
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
	}

	// Create user in database
	if err := database.CreateUser(database.DB, user); err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			http.Error(w, "Email or username already exists", http.StatusConflict)
			return
		}
		log.Printf("Error creating user: %v", err)
		http.Error(w, "Could not create user", http.StatusInternalServerError)
		return
	}

	// Generate JWT verification token
	verificationToken, err := utils.GenerateEmailVerificationToken(user.ID, user.Email)
	if err != nil {
		log.Printf("Error generating verification token: %v", err)
		http.Error(w, "Could not generate verification token", http.StatusInternalServerError)
		return
	}

	// Send verification email
	go func() {
		if err := email.SendVerificationEmail(user.Email, verificationToken); err != nil {
			log.Printf("Warning: Could not send verification email: %v", err)
		}
	}()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":      "success",
		"status_code": http.StatusCreated,
		"message":     "User registered successfully. Please check your email for verification.",
		"user": map[string]interface{}{
			"id":       user.ID,
			"email":    user.Email,
			"username": user.Username,
			"active":   user.Active,
		},
	})
}
