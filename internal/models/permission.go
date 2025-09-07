package models

import "time"

type Permission struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Resource    string    `json:"resource"`
	Action      string    `json:"action"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UserProfile struct {
	ID            string    `json:"id"`
	Username      string    `json:"username"`
	Email         string    `json:"email"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	UserType      string    `json:"user_type"`
	EmailVerified bool      `json:"email_verified"`
	Active        bool      `json:"active"`
	CreatedAt     time.Time `json:"created_at"`
}

type PermissionList struct {
	Permissions []Permission `json:"permissions"`
	Total       int          `json:"total"`
}
