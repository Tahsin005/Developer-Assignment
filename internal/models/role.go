package models

import "time"

type Role struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

type RoleList struct {
	Roles []Role `json:"roles"`
	Total int    `json:"total"`
}
