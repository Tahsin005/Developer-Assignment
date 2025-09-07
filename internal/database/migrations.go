package database

import (
	"fmt"
	"log"
	"os"
)

// Migration represents a database migration
type Migration struct {
	Name         string
	FilePath     string
	SkipIfExists bool
}

// RunMigrations executes all migrations in order
func RunMigrations() error {
	migrations := []Migration{
		{
			Name:     "Initial Schema",
			FilePath: "migrations/000001_init_schema/up.sql",
		},
		{
			Name:         "Add Reset Password Fields",
			FilePath:     "migrations/000002_add_reset_password_fields/up.sql",
			SkipIfExists: true,
		},
	}

	for _, migration := range migrations {
		if err := executeMigration(migration); err != nil {
			return fmt.Errorf("failed to execute migration %s: %v", migration.Name, err)
		}
	}

	return nil
}

// executeMigration runs a single migration
func executeMigration(m Migration) error {
	log.Printf("Running migration: %s", m.Name)

	if m.SkipIfExists {
		// Check if columns exist
		var exists bool
		err := DB.QueryRow(`
			SELECT EXISTS (
				SELECT 1 
				FROM information_schema.columns 
				WHERE table_name = 'users' 
				AND column_name = 'reset_token'
			);
		`).Scan(&exists)

		if err != nil {
			return fmt.Errorf("failed to check column existence: %v", err)
		}

		if exists {
			log.Printf("Skipping migration %s: columns already exist", m.Name)
			return nil
		}
	}

	// Read migration file
	sqlBytes, err := os.ReadFile(m.FilePath)
	if err != nil {
		return fmt.Errorf("failed to read migration file: %v", err)
	}

	// Execute migration
	_, err = DB.Exec(string(sqlBytes))
	if err != nil {
		return fmt.Errorf("failed to execute migration: %v", err)
	}

	log.Printf("Completed migration: %s", m.Name)
	return nil
}
