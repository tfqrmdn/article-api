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

	// Run seeders
	if err := runSeeders(db); err != nil {
		log.Fatal("Failed to run seeders:", err)
	}

	fmt.Println("Seeders completed successfully!")
}

func runSeeders(db *sql.DB) error {
	// Get list of seeder files
	seederFiles, err := getSeederFiles()
	if err != nil {
		return fmt.Errorf("failed to get seeder files: %w", err)
	}

	// Run all seeders
	for _, file := range seederFiles {
		fmt.Printf("Running seeder: %s\n", file)
		if err := runSeeder(db, file); err != nil {
			return fmt.Errorf("failed to run seeder %s: %w", file, err)
		}
	}
	return nil
}

func getSeederFiles() ([]string, error) {
	var files []string

	err := filepath.WalkDir("scripts/seeders", func(path string, d fs.DirEntry, err error) error {
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

func runSeeder(db *sql.DB, filename string) error {
	// Read seeder file
	content, err := os.ReadFile(filepath.Join("scripts/seeders", filename))
	if err != nil {
		return err
	}

	// Execute seeder
	if _, err := db.Exec(string(content)); err != nil {
		return err
	}

	return nil
}
