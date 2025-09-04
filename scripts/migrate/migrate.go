package main

import (
	"database/sql"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"article-api/internal/database"

	_ "github.com/lib/pq"
)

func main() {
	// Connect to database
	db, err := database.Connect()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Run migrations
	if err := runMigrations(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	fmt.Println("Migrations completed successfully!")
}

func runMigrations(db *sql.DB) error {
	// Create migrations table if it doesn't exist
	createMigrationsTable := `
		CREATE TABLE IF NOT EXISTS migrations (
			id SERIAL PRIMARY KEY,
			filename TEXT NOT NULL UNIQUE,
			executed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`
	if _, err := db.Exec(createMigrationsTable); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get list of migration files
	migrationFiles, err := getMigrationFiles()
	if err != nil {
		return fmt.Errorf("failed to get migration files: %w", err)
	}

	// Get executed migrations
	executedMigrations, err := getExecutedMigrations(db)
	if err != nil {
		return fmt.Errorf("failed to get executed migrations: %w", err)
	}

	// Run pending migrations
	for _, file := range migrationFiles {
		if !executedMigrations[file] {
			fmt.Printf("Running migration: %s\n", file)
			if err := runMigration(db, file); err != nil {
				return fmt.Errorf("failed to run migration %s: %w", file, err)
			}
		} else {
			fmt.Printf("Skipping already executed migration: %s\n", file)
		}
	}

	return nil
}

func getMigrationFiles() ([]string, error) {
	var files []string

	err := filepath.WalkDir("scripts/migrations", func(path string, d fs.DirEntry, err error) error {
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

func getExecutedMigrations(db *sql.DB) (map[string]bool, error) {
	rows, err := db.Query("SELECT filename FROM migrations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	executed := make(map[string]bool)
	for rows.Next() {
		var filename string
		if err := rows.Scan(&filename); err != nil {
			return nil, err
		}
		executed[filename] = true
	}

	return executed, nil
}

func runMigration(db *sql.DB, filename string) error {
	// Read migration file
	content, err := os.ReadFile(filepath.Join("scripts/migrations", filename))
	if err != nil {
		return err
	}

	// Execute migration
	if _, err := db.Exec(string(content)); err != nil {
		return err
	}

	// Record migration as executed
	_, err = db.Exec("INSERT INTO migrations (filename) VALUES ($1)", filename)
	return err
}
