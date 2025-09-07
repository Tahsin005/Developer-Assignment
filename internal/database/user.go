package database

import (
	"database/sql"
	"fmt"

	"github.com/soyaibzihad10/Developer-Assignment/internal/models"
)

var rolePriority = map[string]int{
	"user":         1,
	"moderator":    2,
	"admin":        3,
	"system_admin": 4,
}


func ListUsers() ([]models.User, error) {
	rows, err := DB.Query(`
        SELECT id, username, email, first_name, last_name, 
               email_verified, user_type, active, created_at
        FROM users 
        WHERE deletion_requested = false
        ORDER BY created_at DESC
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		err := rows.Scan(
			&u.ID, &u.Username, &u.Email, &u.FirstName, &u.LastName,
			&u.EmailVerified, &u.UserType, &u.Active, &u.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func UpdateUser(db *sql.DB, userID string, req *models.UserUpdateRequest) (*models.User, error) {
	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `
        UPDATE users 
        SET first_name = COALESCE($1, first_name),
            last_name = COALESCE($2, last_name),
            email = COALESCE($3, email),
            username = COALESCE($4, username),
            updated_at = NOW()
        WHERE id = $5 AND deletion_requested = false
        RETURNING id, username, email, first_name, last_name, user_type, 
                  email_verified, active, created_at, updated_at
    `

	user := &models.User{}
	err = tx.QueryRow(
		query,
		req.FirstName,
		req.LastName,
		req.Email,
		req.Username,
		userID,
	).Scan(
		&user.ID, &user.Username, &user.Email, &user.FirstName, &user.LastName,
		&user.UserType, &user.EmailVerified, &user.Active, &user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return user, nil
}

func DeleteUser(db *sql.DB, userID string) error {
	result, err := db.Exec("DELETE FROM users WHERE id = $1", userID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func RequestUserDeletion(db *sql.DB, userID string) error {
	result, err := db.Exec(`
        UPDATE users 
        SET deletion_requested = true,
            updated_at = NOW()
        WHERE id = $1 AND deletion_requested = false
    `, userID)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}


// Promote to admin only if user is below admin
func PromoteUserToAdmin(db *sql.DB, userID string) error {
	currentRole, err := getCurrentUserRole(db, userID)
	if err != nil {
		return err
	}
	if rolePriority[currentRole] >= rolePriority["admin"] {
		return fmt.Errorf("user already has role '%s' or higher", currentRole)
	}
	return ChangeUserRole(db, userID, "admin")
}

// Promote to moderator only if user is below moderator
func PromoteUserToModerator(db *sql.DB, userID string) error {
	currentRole, err := getCurrentUserRole(db, userID)
	if err != nil {
		return err
	}
	if rolePriority[currentRole] >= rolePriority["moderator"] {
		return fmt.Errorf("user already has role '%s' or higher", currentRole)
	}
	return ChangeUserRole(db, userID, "moderator")
}

// Demote based on hierarchy
func DemoteUser(db *sql.DB, userID string) error {
	currentRole, err := getCurrentUserRole(db, userID)
	if err != nil {
		return err
	}

	var newRole string
	switch currentRole {
	case "system_admin":
		newRole = "admin"
	case "admin":
		newRole = "moderator"
	case "moderator":
		newRole = "user"
	default:
		return fmt.Errorf("cannot demote user with role: %s", currentRole)
	}

	return ChangeUserRole(db, userID, newRole)
}

// Helper to fetch current role
func getCurrentUserRole(db *sql.DB, userID string) (string, error) {
	var role string
	err := db.QueryRow("SELECT user_type FROM users WHERE id = $1", userID).Scan(&role)
	return role, err
}

// Change role in both `users` and `user_roles` table
func ChangeUserRole(db *sql.DB, userID string, newRole string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
        UPDATE users 
        SET user_type = $1,
            updated_at = NOW()
        WHERE id = $2
    `, newRole, userID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
        UPDATE user_roles
        SET role_id = (SELECT id FROM roles WHERE name = $1)
        WHERE user_id = $2
    `, newRole, userID)
	if err != nil {
		return err
	}

	return tx.Commit()
}