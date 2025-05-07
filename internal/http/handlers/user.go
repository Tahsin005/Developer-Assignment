package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/tahsin005/affpilot-auth/internal/database"
	"github.com/tahsin005/affpilot-auth/internal/http/middleware"
	"github.com/tahsin005/affpilot-auth/internal/models"
	"github.com/tahsin005/affpilot-auth/internal/utils"
)

func UsersListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	queryStatement := `
		SELECT id, username, email, first_name, last_name, email_verified, user_type, 
		verification_token, token_expiry, deletion_requested, active, created_at, updated_at
		FROM users
	`

	rows, err := database.DB.Query(queryStatement)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Database query error")
		log.Println("Query error:", err)
		return
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User
		var tokenExpiry sql.NullTime
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.EmailVerified,
			&user.UserType,
			&user.VerificationToken,
			&tokenExpiry,
			&user.DeletionRequested,
			&user.Active,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "Failed to scan user")
			log.Println("Scan error:", err)
			return
		}

		if tokenExpiry.Valid {
			user.TokenExpiry = &tokenExpiry.Time
		}

		users = append(users, user)
	}


	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User's List",
		"users": users,
	})
}


func UserDetailsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	userIDStr := vars["user_id"]

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid user_id format")
		return
	}

	query := `
		SELECT id, username, email, first_name, last_name, email_verified, user_type, 
		verification_token, token_expiry, deletion_requested, active, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user models.User
	var tokenExpiry sql.NullTime

	err = database.DB.QueryRow(query, userID).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.EmailVerified,
		&user.UserType,
		&user.VerificationToken,
		&tokenExpiry,
		&user.DeletionRequested,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			utils.WriteError(w, http.StatusNotFound, "User not found")
		} else {
			utils.WriteError(w, http.StatusInternalServerError, "Database query error")
			log.Println("Query error:", err)
		}
		return
	}

	if tokenExpiry.Valid {
		user.TokenExpiry = &tokenExpiry.Time
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User details fetched successfully",
		"user":    user,
	})
}


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

	if req.Username == "" || req.Email == "" || req.FirstName == "" || req.LastName == "" {
		utils.WriteError(w, http.StatusBadRequest, "All fields are required")
		return
	}

	if !utils.IsValidEmail(req.Email) {
		utils.WriteError(w, http.StatusBadRequest, "Invalid email format")
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

	emailVerified := user.EmailVerified
    if req.Email != user.Email {
        emailVerified = false
    }

	queryStatement = `
        UPDATE users
        SET username = $1, email = $2, first_name = $3, last_name = $4, 
            email_verified = $5, updated_at = NOW()
        WHERE id = $6
        RETURNING updated_at
    `

	var updatedAt time.Time
	err = database.DB.QueryRow(queryStatement, req.Username, req.Email, req.FirstName, req.LastName, emailVerified, userID).Scan(&updatedAt)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to update user")
		return
	}

	user.Username = req.Username
    user.Email = req.Email
    user.FirstName = req.FirstName
    user.LastName = req.LastName
    user.EmailVerified = emailVerified
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


func UserDeletionRequestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	claims, ok := r.Context().Value("user").(*middleware.Claims)
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	vars := mux.Vars(r)
	userIDStr := vars["user_id"]

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid user_id format")
		return
	}

	if userIDStr != claims.UserID {
		utils.WriteError(w, http.StatusForbidden, "Forbidden: You can only request deletion for your own account")
		return
	}

	user, err := utils.GetUserByID(userID)
	if err == sql.ErrNoRows {
		utils.WriteError(w, http.StatusNotFound, "User not found")
		return
	}
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	if user.DeletionRequested {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "You have already requested deletion; please wait for the update",
		})
		return
	}

	queryStatement := `
		UPDATE users
		SET deletion_requested = TRUE, updated_at = $1
		WHERE id = $2 AND active = TRUE
	`
	result, err := database.DB.Exec(queryStatement, time.Now(), userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to request account deletion")
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	if rowsAffected == 0 {
		utils.WriteError(w, http.StatusNotFound, "User not found or already inactive")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Your deletion request has been sent successfully",
	})
}


func UserDeleteHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	claims, ok := r.Context().Value("user").(*middleware.Claims)
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	vars := mux.Vars(r)
	userIDStr := vars["user_id"]

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid user_id format")
		return
	}

	if userIDStr == claims.UserID {
		utils.WriteError(w, http.StatusForbidden, "Forbidden: You cannot delete your own account")
		return
	}

	user, err := utils.GetUserByID(userID)
	if err == sql.ErrNoRows {
		utils.WriteError(w, http.StatusNotFound, "User not found or already deleted")
		return
	}
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// preventing deletion of system admins
	if user.UserType == "system_admin" {
		utils.WriteError(w, http.StatusForbidden, "Forbidden: System admins cannot be deleted")
		return
	}

	switch claims.Role {
	case "moderator":
		// moderators can only delete users who requested deletion and are not admins
		if !user.DeletionRequested {
			utils.WriteError(w, http.StatusForbidden, "Forbidden: User has not requested deletion")
			return
		}
		if user.UserType == "admin" {
			utils.WriteError(w, http.StatusForbidden, "Forbidden: Moderators cannot delete admins")
			return
		}
	case "admin":
		// admins cannot delete other admins
		if user.UserType == "admin" {
			utils.WriteError(w, http.StatusForbidden, "Forbidden: Admins cannot delete other admins")
			return
		}
	case "system_admin":
		// system admins can delete anyone except system admins (already checked)
		log.Printf("System admin %s authorized to delete user %s", claims.UserID, userIDStr)
	default:
		utils.WriteError(w, http.StatusForbidden, "Forbidden: Insufficient role")
		return
	}

	// transaction
	tx, err := database.DB.Begin()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	defer tx.Rollback()

	// soft deletion
	queryStatement := `
		UPDATE users
		SET active = FALSE, updated_at = $1
		WHERE id = $2 AND active = TRUE
	`
	result, err := tx.Exec(queryStatement, time.Now(), userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to delete user")
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	if rowsAffected == 0 {
		utils.WriteError(w, http.StatusNotFound, "User not found or already deleted")
		return
	}

	queryStatement = `
		DELETE FROM user_roles
		WHERE user_id = $1
	`
	_, err = tx.Exec(queryStatement, userID)
	if err != nil {
		log.Printf("Warning: Failed to clean up user_roles for user %s: %v", userID, err)
	}

	if err := tx.Commit(); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to complete deletion")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User deleted successfully",
	})
}


func UserRoleChangeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	claims, ok := r.Context().Value("user").(*middleware.Claims)
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	vars := mux.Vars(r)
	userIDStr := vars["user_id"]

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid user_id format")
		return
	}

	if userIDStr == claims.UserID {
		utils.WriteError(w, http.StatusForbidden, "Forbidden: You cannot change your own role")
		return
	}

	var req models.RoleChangeRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	req.Role = strings.ToLower(strings.TrimSpace(req.Role))
	validRoles := map[string]bool{
		"user": true,
		"moderator": true,
		"admin": true,
		"system_admin": true,
	}

	if !validRoles[req.Role] {
		utils.WriteError(w, http.StatusBadRequest, "Invalid role: only accepts user, moderator, admin & system_admin")
		return
	}

	user, err := utils.GetUserByID(userID)
	if err == sql.ErrNoRows {
		utils.WriteError(w, http.StatusNotFound, "User not found or already deleted")
		return
	}

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return	
	}

	// preventing role change for system admin
	if user.UserType == "system_admin" {
		utils.WriteError(w, http.StatusForbidden, "Forbidden: System admins' roles cannot be changed")
		return
	}

	// role change rules
	switch claims.Role {
	case "admin":
		if req.Role == "system_admin" {
			utils.WriteError(w, http.StatusForbidden, "Forbidden: Admins cannot assign system_admin role")
			return
		}
		if user.UserType == "admin" {
			utils.WriteError(w, http.StatusForbidden, "Forbidden: Admins cannot change other admins' roles")
			return
		}
		if user.UserType == "moderator" && req.Role != "user" && req.Role != "admin" {
			utils.WriteError(w, http.StatusForbidden, "Forbidden: Moderators can only be changed to user or admin")
			return
		}
		if user.UserType == "user" && req.Role != "moderator" && req.Role != "admin" {
			utils.WriteError(w, http.StatusForbidden, "Forbidden: Users can only be changed to moderator or admin")
			return
		}
	case "system_admin":
		log.Printf("System admin %s authorized to change role of user %s to %s", claims.UserID, userIDStr, req.Role)
	default:
		utils.WriteError(w, http.StatusForbidden, "Forbidden: Insufficient role")
		return
	}


	// transaction
	tx, err := database.DB.Begin()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	defer tx.Rollback()

	queryStatement := `
		SELECT role_id 
		FROM user_roles
		WHERE user_id = $1
	`
	var existingRoleID uuid.UUID
	err = tx.QueryRow(queryStatement, userID).Scan(&existingRoleID)

	if err == sql.ErrNoRows {
		query := `
			INSERT INTO user_roles (user_id, role_id, assigned_by, created_at)
			VALUES ($1, (SELECT id FROM roles WHERE name = $2), $3, $4)
		`
		_, err = tx.Exec(query, userID, req.Role, claims.UserID, time.Now())
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "Failed to assign new role")
			return
		}
	} else if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to check existing role")
		return
	} else {
		query := `
			UPDATE user_roles
			SET role_id = (SELECT id FROM roles WHERE name = $1), assigned_by = $2, created_at = $3
			WHERE user_id = $4
		`
		result, err := tx.Exec(query, req.Role, claims.UserID, time.Now(), userID)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "Failed to update user role")
			return
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		if rowsAffected == 0 {
			utils.WriteError(w, http.StatusInternalServerError, "Failed to update user role")
			return
		}
	}

	query := `
		UPDATE users
		SET user_type = $1, updated_at = $2
		WHERE id = $3 AND active = TRUE
	`
	result, err := tx.Exec(query, req.Role, time.Now(), userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to update user type")
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	if rowsAffected == 0 {
		utils.WriteError(w, http.StatusNotFound, "User not found or already deleted")
		return
	}

	if err := tx.Commit(); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to complete role change")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("User role changed to %s successfully", req.Role),
	})
}