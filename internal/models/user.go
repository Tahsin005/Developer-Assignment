package models

type User struct {
    ID                string `json:"id"`
    Username          string `json:"username"`
    Email             string `json:"email"`
    PasswordHash      string `json:"password_hash"`
    FirstName         string `json:"first_name,omitempty"`
    LastName          string `json:"last_name,omitempty"`
    EmailVerified     bool   `json:"email_verified"`
    UserType          string `json:"user_type"`
    VerificationToken string `json:"verification_token,omitempty"`
    TokenExpiry       string `json:"token_expiry,omitempty"`
    DeletionRequested bool   `json:"deletion_requested"`
    Active            bool   `json:"active"`
    CreatedAt         string `json:"created_at"`
    UpdatedAt         string `json:"updated_at"`
}

type RegisterRequest struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}