package database

import (
	"github.com/soyaibzihad10/Developer-Assignment/internal/models"
)

// ListPermissions returns all permissions
func ListPermissions() ([]models.Permission, error) {
	rows, err := DB.Query(`
        SELECT id, name, resource, action, description, created_at, updated_at 
        FROM permissions 
        ORDER BY name
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []models.Permission
	for rows.Next() {
		var p models.Permission
		err := rows.Scan(
			&p.ID, &p.Name, &p.Resource, &p.Action,
			&p.Description, &p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, p)
	}
	return permissions, nil
}

// GetPermissionByID returns a single permission
func GetPermissionByID(id string) (*models.Permission, error) {
	p := &models.Permission{}
	err := DB.QueryRow(`
        SELECT id, name, resource, action, description, created_at, updated_at 
        FROM permissions 
        WHERE id = $1
    `, id).Scan(
		&p.ID, &p.Name, &p.Resource, &p.Action,
		&p.Description, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return p, nil
}

// GetUserPermissions returns all permissions for a user
func GetUserPermissions(userID string) ([]models.Permission, error) {
	rows, err := DB.Query(`
        SELECT DISTINCT
            p.id,
            p.name,
            p.resource,
            p.action,
            p.description,
            p.created_at,
            p.updated_at
        FROM permissions p
        INNER JOIN role_permissions rp ON p.id = rp.permission_id
        INNER JOIN user_roles ur ON rp.role_id = ur.role_id
        WHERE ur.user_id = $1
        ORDER BY p.name`, userID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []models.Permission
	for rows.Next() {
		var p models.Permission
		err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Resource,
			&p.Action,
			&p.Description,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, p)
	}

	return permissions, nil
}
