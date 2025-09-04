package migration

import (
	"database/sql"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// RunMigrations runs all pending migrations
func RunMigrations(db *sql.DB) error {
	log.Println("Starting database migrations...")

	// Create migrations table if it doesn't exist
	if err := createMigrationsTable(db); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get list of migration files
	migrationFiles, err := getMigrationFiles()
	if err != nil {
		return fmt.Errorf("failed to get migration files: %w", err)
	}

	// Get list of executed migrations
	executedMigrations, err := getExecutedMigrations(db)
	if err != nil {
		return fmt.Errorf("failed to get executed migrations: %w", err)
	}

	// Execute pending migrations
	executedCount := 0
	for _, file := range migrationFiles {
		if !isMigrationExecuted(file, executedMigrations) {
			if err := executeMigration(db, file); err != nil {
				return fmt.Errorf("failed to execute migration %s: %w", file, err)
			}
			executedCount++
		} else {
			log.Printf("Skipping already executed migration: %s", file)
		}
	}

	if executedCount > 0 {
		log.Printf("Successfully executed %d migrations", executedCount)
	} else {
		log.Println("No pending migrations found")
	}

	log.Println("Database migrations completed successfully!")
	return nil
}

// createMigrationsTable creates the migrations tracking table
func createMigrationsTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS migrations (
			id SERIAL PRIMARY KEY,
			filename VARCHAR(255) NOT NULL UNIQUE,
			executed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`
	_, err := db.Exec(query)
	return err
}

// getMigrationFiles returns a sorted list of migration files
func getMigrationFiles() ([]string, error) {
	var files []string

	// Look for migration files in the scripts/migrations directory
	migrationDir := "scripts/migrations"

	// Check if running in Docker (different path)
	if _, err := os.Stat(migrationDir); os.IsNotExist(err) {
		migrationDir = "/app/scripts/migrations"
	}

	err := filepath.WalkDir(migrationDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && strings.HasSuffix(path, ".sql") {
			files = append(files, filepath.Base(path))
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	sort.Strings(files)
	return files, nil
}

// getExecutedMigrations returns a list of already executed migration filenames
func getExecutedMigrations(db *sql.DB) ([]string, error) {
	query := "SELECT filename FROM migrations ORDER BY filename"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var migrations []string
	for rows.Next() {
		var filename string
		if err := rows.Scan(&filename); err != nil {
			return nil, err
		}
		migrations = append(migrations, filename)
	}

	return migrations, rows.Err()
}

// isMigrationExecuted checks if a migration has already been executed
func isMigrationExecuted(filename string, executedMigrations []string) bool {
	for _, executed := range executedMigrations {
		if executed == filename {
			return true
		}
	}
	return false
}

// executeMigration executes a single migration file
func executeMigration(db *sql.DB, filename string) error {
	// Determine the correct path
	migrationDir := "scripts/migrations"
	if _, err := os.Stat(migrationDir); os.IsNotExist(err) {
		migrationDir = "/app/scripts/migrations"
	}

	filePath := filepath.Join(migrationDir, filename)

	// Read migration file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	// Execute migration
	_, err = db.Exec(string(content))
	if err != nil {
		return fmt.Errorf("failed to execute migration SQL: %w", err)
	}

	// Record migration as executed
	_, err = db.Exec("INSERT INTO migrations (filename) VALUES ($1)", filename)
	if err != nil {
		return fmt.Errorf("failed to record migration: %w", err)
	}

	log.Printf("Successfully executed migration: %s", filename)
	return nil
}
