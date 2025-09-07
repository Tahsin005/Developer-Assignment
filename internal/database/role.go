package database

import (
	"time"

	"github.com/google/uuid"
	"github.com/soyaibzihad10/Developer-Assignment/internal/models"
)

import "errors"

var (
    ErrNotFound = errors.New("record not found")
)

func ListRoles() ([]models.Role, error) {
	rows, err := DB.Query(`
        SELECT id, name, description, created_at, updated_at 
        FROM roles 
        ORDER BY created_at DESC
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []models.Role
	for rows.Next() {
		var role models.Role
		err := rows.Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, nil
}

func GetRoleByID(id string) (*models.Role, error) {
	role := &models.Role{}
	err := DB.QueryRow(`
        SELECT id, name, description, created_at, updated_at 
        FROM roles 
        WHERE id = $1
    `, id).Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return role, nil
}

func CreateRole(role *models.Role) error {
	role.ID = uuid.New().String()
	now := time.Now()
	role.CreatedAt = now
	role.UpdatedAt = now

	_, err := DB.Exec(`
        INSERT INTO roles (id, name, description, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5)
    `, role.ID, role.Name, role.Description, role.CreatedAt, role.UpdatedAt)

	return err
}

func UpdateRole(id string, role *models.Role) error {
	role.UpdatedAt = time.Now()

	result, err := DB.Exec(`
        UPDATE roles 
        SET name = $1, description = $2, updated_at = $3
        WHERE id = $4
    `, role.Name, role.Description, role.UpdatedAt, id)

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

func DeleteRole(id string) error {
	result, err := DB.Exec("DELETE FROM roles WHERE id = $1", id)
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
