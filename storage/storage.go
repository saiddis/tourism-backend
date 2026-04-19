package storage

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"tourism-backend/config"

	_ "github.com/lib/pq"
)

func NewPostgresDB(cfg *config.Config) (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening db: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to db: %w", err)
	}

	return db, nil
}

func RunMigrations(db *sql.DB, migrationsDir string) error {
	// Check if migrations table exists
	var exists bool
	err := db.QueryRow(`
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = 'schema_migrations'
		)
	`).Scan(&exists)

	if err != nil {
		return fmt.Errorf("error checking migrations table: %w", err)
	}

	// Create migrations tracking table if it doesn't exist
	if !exists {
		_, err = db.Exec(`
			CREATE TABLE IF NOT EXISTS schema_migrations (
				id SERIAL PRIMARY KEY,
				filename VARCHAR(255) NOT NULL UNIQUE,
				applied_at TIMESTAMP NOT NULL DEFAULT NOW()
			)
		`)
		if err != nil {
			return fmt.Errorf("error creating migrations table: %w", err)
		}
	}

	// Read migration files
	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("error reading migrations directory: %w", err)
	}

	// Apply each migration file
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".sql" {
			continue
		}

		filename := entry.Name()

		// Check if already applied
		var applied bool
		err := db.QueryRow(`
			SELECT EXISTS (
				SELECT FROM schema_migrations 
				WHERE filename = $1
			)
		`, filename).Scan(&applied)

		if err != nil {
			return fmt.Errorf("error checking migration %s: %w", filename, err)
		}

		if applied {
			fmt.Printf("Migration %s already applied, skipping\n", filename)
			continue
		}

		// Read and apply migration
		content, err := os.ReadFile(filepath.Join(migrationsDir, filename))
		if err != nil {
			return fmt.Errorf("error reading migration %s: %w", filename, err)
		}

		_, err = db.Exec(string(content))
		if err != nil {
			return fmt.Errorf("error applying migration %s: %w", filename, err)
		}

		// Record migration as applied
		_, err = db.Exec(`
			INSERT INTO schema_migrations (filename) VALUES ($1)
		`, filename)
		if err != nil {
			return fmt.Errorf("error recording migration %s: %w", filename, err)
		}

		fmt.Printf("Applied migration: %s\n", filename)
	}

	return nil
}
