package database

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// RunMigrations reads SQL files from the given directory and executes any that
// have not yet been applied. Applied migrations are tracked in the
// schema_migrations table.
func (db *DB) RunMigrations(migrationsDir string) error {
	// Ensure the tracking table exists
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version TEXT PRIMARY KEY,
			applied_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)
	`); err != nil {
		return fmt.Errorf("creating schema_migrations table: %w", err)
	}

	// Collect already-applied versions
	applied, err := db.appliedMigrations()
	if err != nil {
		return err
	}

	// Read migration files
	files, err := filepath.Glob(filepath.Join(migrationsDir, "*.sql"))
	if err != nil {
		return fmt.Errorf("reading migration files: %w", err)
	}
	sort.Strings(files)

	for _, file := range files {
		version := strings.TrimSuffix(filepath.Base(file), ".sql")
		if applied[version] {
			continue
		}

		log.Printf("applying migration: %s", version)

		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("reading migration %s: %w", version, err)
		}

		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("starting transaction for %s: %w", version, err)
		}

		if _, err := tx.Exec(string(content)); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("executing migration %s: %w", version, err)
		}

		if _, err := tx.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", version); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("recording migration %s: %w", version, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("committing migration %s: %w", version, err)
		}

		log.Printf("migration applied: %s", version)
	}

	return nil
}

func (db *DB) appliedMigrations() (map[string]bool, error) {
	applied := make(map[string]bool)

	rows, err := db.Query("SELECT version FROM schema_migrations ORDER BY version")
	if err != nil {
		return nil, fmt.Errorf("querying applied migrations: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, fmt.Errorf("scanning migration version: %w", err)
		}
		applied[version] = true
	}

	return applied, rows.Err()
}
