// Package database manages the PostgreSQL connection pool and health checks.
package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

// DB wraps the standard sql.DB connection pool.
type DB struct {
	*sql.DB
}

// Connect establishes a connection pool to PostgreSQL and verifies connectivity.
func Connect(databaseURL string) (*DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("opening database connection: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(1 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("pinging database: %w", err)
	}

	return &DB{db}, nil
}

// HealthCheck verifies the database connection is alive.
func (db *DB) HealthCheck() error {
	return db.Ping()
}
