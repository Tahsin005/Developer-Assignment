package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/soyaibzihad10/Developer-Assignment/internal/models"
)

var ErrUserNotFound = errors.New("user not found")

func UpdateEmailVerificationStatus(db *sql.DB, userID string) error {
	updateQuery := `
        UPDATE users 
        SET email_verified = true
        WHERE id = $1 AND email_verified = false
    `

	result, err := db.Exec(updateQuery, userID)
	if err != nil {
		return fmt.Errorf("failed to update verification status: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no user found with ID: %s", userID)
	}

	return nil
}

func UpdateUserPassword(db *sql.DB, userID string, hashedPassword string) error {
	query := `UPDATE users SET 
        password_hash = $1,
        reset_token = NULL,
        reset_token_expiry = NULL
        WHERE id = $2`
	_, err := db.Exec(query, hashedPassword, userID)
	return err
}

func UpdateUserVerificationToken(db *sql.DB, user *models.User) error {
	query := "UPDATE users SET verification_token = $1, token_expiry = $2 WHERE id = $3"
	_, err := db.Exec(query, user.VerificationToken, user.TokenExpiry, user.ID)
	return err
}

func CreateUser(db *sql.DB, user *models.User) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO users (
        username, email, first_name, last_name, password_hash,
        email_verified, user_type, active, created_at, updated_at
    ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
    RETURNING id`

	now := time.Now()

	// Set default values
	user.CreatedAt = now
	user.UpdatedAt = now
	user.EmailVerified = false
	user.UserType = "user"
	user.Active = true

	err = tx.QueryRow(query,
		user.Username,
		user.Email,
		user.FirstName,
		user.LastName,
		user.PasswordHash,
		user.EmailVerified,
		user.UserType,
		user.Active,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.ID)

	if err != nil {
		return err
	}

	// Assign default user role
	roleQuery := `
        INSERT INTO user_roles (user_id, role_id, assigned_by)
        SELECT $1, id, $1
        FROM roles
        WHERE name = 'user'
    `
	_, err = tx.Exec(roleQuery, user.ID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func GetUserByID(db *sql.DB, id string) (*models.User, error) {
	user := &models.User{}
	err := db.QueryRow(`
        SELECT id, username, email, first_name, last_name, 
               user_type, email_verified, active, created_at, updated_at
        FROM users 
        WHERE id = $1 AND deletion_requested = false
    `, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.FirstName, &user.LastName,
		&user.UserType, &user.EmailVerified, &user.Active, &user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func GetUserByEmail(db *sql.DB, email string) (*models.User, error) {
	user := &models.User{}
	query := "SELECT id, username, first_name, last_name, email, password_hash, email_verified, user_type, active, created_at, updated_at FROM users WHERE email = $1"
	err := db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Username,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.PasswordHash,
		&user.EmailVerified,
		&user.UserType,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}
