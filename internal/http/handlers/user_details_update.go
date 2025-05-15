package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/lib/pq"
	"github.com/tahsin005/affpilot-auth/internal/database"
	"github.com/tahsin005/affpilot-auth/internal/models"
	"github.com/tahsin005/affpilot-auth/internal/utils"
)

func UserUpdateDetailsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	userIDStr := vars["user_id"]

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid user_id format")
		return
	}

	var req models.UserUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid input")
		return
	}

	if req.Username == "" || req.FirstName == "" || req.LastName == "" {
		utils.WriteError(w, http.StatusBadRequest, "Username, first name, and last name are required")
		return
	}

	queryStatement := `
		SELECT id, username, email, first_name, last_name, email_verified, user_type, 
		verification_token, token_expiry, deletion_requested, active, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user models.User
	err = database.DB.QueryRow(queryStatement, userID).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.EmailVerified,
		&user.UserType,
		&user.VerificationToken,
		&user.TokenExpiry,
		&user.DeletionRequested,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		utils.WriteError(w, http.StatusNotFound, "User not found")
		return
	}
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to retrieve user")
		return
	}

	queryStatement = `
        UPDATE users
        SET username = $1, first_name = $2, last_name = $3, updated_at = NOW()
        WHERE id = $4
        RETURNING updated_at
    `

	var updatedAt time.Time
	err = database.DB.QueryRow(queryStatement, req.Username, req.FirstName, req.LastName, userID).Scan(&updatedAt)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			utils.WriteError(w, http.StatusBadRequest, "Username is already taken. Please choose another.")
			return
		}

		utils.WriteError(w, http.StatusInternalServerError, "Failed to update user")
		return
	}


	user.Username = req.Username
	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.UpdatedAt = updatedAt

	response := models.UpdateUserResponse{
		ID:            user.ID,
		Username:      user.Username,
		Email:         user.Email,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		EmailVerified: user.EmailVerified,
		UserType:      user.UserType,
		Active:        user.Active,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User details updated successfully",
		"user":    response,
	})
}
