// Package clients manages client (horse owner) data access and business logic.
package clients

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/stride-pro/backend/internal/database"
	"github.com/stride-pro/backend/internal/models"
)

// Repository handles client persistence.
type Repository struct {
	db *database.DB
}

// NewRepository creates a client repository.
func NewRepository(db *database.DB) *Repository {
	return &Repository{db: db}
}

// Create inserts a new client.
func (r *Repository) Create(c *models.Client) error {
	c.ID = uuid.New()
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()

	_, err := r.db.Exec(`
		INSERT INTO clients (id, user_id, first_name, last_name, email, phone, address, notes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		c.ID, c.UserID, c.FirstName, c.LastName, c.Email, c.Phone, c.Address, c.Notes, c.CreatedAt, c.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("inserting client: %w", err)
	}
	return nil
}

// GetByID returns a single client by ID, scoped to the user, including associated horses.
func (r *Repository) GetByID(userID, clientID uuid.UUID) (*models.Client, error) {
	c := &models.Client{}
	err := r.db.QueryRow(`
		SELECT id, user_id, first_name, last_name, email, phone, address, notes, created_at, updated_at
		FROM clients WHERE id = $1 AND user_id = $2`,
		clientID, userID,
	).Scan(&c.ID, &c.UserID, &c.FirstName, &c.LastName, &c.Email, &c.Phone, &c.Address, &c.Notes, &c.CreatedAt, &c.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("querying client: %w", err)
	}

	rows, err := r.db.Query(`
		SELECT id, name, breed, age, gender
		FROM horses
		WHERE user_id = $1 AND client_id = $2
		ORDER BY name`, userID, clientID,
	)
	if err != nil {
		return nil, fmt.Errorf("querying client horses: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		h := &models.Horse{}
		if err := rows.Scan(&h.ID, &h.Name, &h.Breed, &h.Age, &h.Gender); err != nil {
			return nil, fmt.Errorf("scanning client horse: %w", err)
		}
		c.Horses = append(c.Horses, h)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating client horses: %w", err)
	}
	return c, nil
}

// GetAllByUserID returns all clients belonging to a user.
func (r *Repository) GetAllByUserID(userID uuid.UUID) ([]models.Client, error) {
	rows, err := r.db.Query(`
		SELECT id, user_id, first_name, last_name, email, phone, address, notes, created_at, updated_at
		FROM clients WHERE user_id = $1 ORDER BY last_name, first_name`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("querying clients: %w", err)
	}
	defer rows.Close()

	var clients []models.Client
	for rows.Next() {
		var c models.Client
		if err := rows.Scan(&c.ID, &c.UserID, &c.FirstName, &c.LastName, &c.Email, &c.Phone, &c.Address, &c.Notes, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scanning client: %w", err)
		}
		clients = append(clients, c)
	}
	return clients, rows.Err()
}

// Update modifies an existing client.
func (r *Repository) Update(c *models.Client) error {
	c.UpdatedAt = time.Now()
	result, err := r.db.Exec(`
		UPDATE clients SET first_name=$1, last_name=$2, email=$3, phone=$4, address=$5, notes=$6, updated_at=$7
		WHERE id=$8 AND user_id=$9`,
		c.FirstName, c.LastName, c.Email, c.Phone, c.Address, c.Notes, c.UpdatedAt, c.ID, c.UserID,
	)
	if err != nil {
		return fmt.Errorf("updating client: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("client not found")
	}
	return nil
}

// CountByUserID returns the number of clients belonging to a user.
func (r *Repository) CountByUserID(userID uuid.UUID) (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM clients WHERE user_id = $1", userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("counting clients: %w", err)
	}
	return count, nil
}

// Delete removes a client by ID, scoped to the user.
func (r *Repository) Delete(userID, clientID uuid.UUID) error {
	result, err := r.db.Exec("DELETE FROM clients WHERE id = $1 AND user_id = $2", clientID, userID)
	if err != nil {
		return fmt.Errorf("deleting client: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("client not found")
	}
	return nil
}
