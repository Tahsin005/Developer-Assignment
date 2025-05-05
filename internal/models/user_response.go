package models

type UserResponse struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	UserType  string `json:"user_type"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
