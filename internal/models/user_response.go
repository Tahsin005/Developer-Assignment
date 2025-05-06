package models

import (
	"time"

	"github.com/google/uuid"
)

type RegisterUserResponse struct {
	ID        uuid.UUID `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	UserType  string `json:"user_type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UpdateUserResponse struct {
    ID            uuid.UUID `json:"id"`
    Username      string    `json:"username"`
    Email         string    `json:"email"`
    FirstName     string    `json:"first_name"`
    LastName      string    `json:"last_name"`
    EmailVerified bool      `json:"email_verified"`
    UserType      string    `json:"user_type"`
    Active        bool      `json:"active"`
    CreatedAt     time.Time `json:"created_at"`
    UpdatedAt     time.Time `json:"updated_at"`
}