// Package service_items manages the user's reusable service price catalog.
package service_items

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/stride-pro/backend/internal/database"
	"github.com/stride-pro/backend/internal/models"
)

// Repository handles service item persistence.
type Repository struct {
	db *database.DB
}

// NewRepository creates a service item repository.
func NewRepository(db *database.DB) *Repository {
	return &Repository{db: db}
}

// GetAll returns all service items for a user.
func (r *Repository) GetAll(userID uuid.UUID) ([]models.ServiceItem, error) {
	rows, err := r.db.Query(`
		SELECT id, user_id, name, default_price, created_at, updated_at
		FROM service_items WHERE user_id = $1 ORDER BY name ASC`, userID,
	)
	if err != nil {
		return nil, fmt.Errorf("querying service items: %w", err)
	}
	defer rows.Close()

	var items []models.ServiceItem
	for rows.Next() {
		var item models.ServiceItem
		if err := rows.Scan(&item.ID, &item.UserID, &item.Name, &item.DefaultPrice, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scanning service item: %w", err)
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

// Create inserts a new service item.
func (r *Repository) Create(item *models.ServiceItem) error {
	item.ID = uuid.New()
	item.CreatedAt = time.Now()
	item.UpdatedAt = time.Now()

	_, err := r.db.Exec(`
		INSERT INTO service_items (id, user_id, name, default_price, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		item.ID, item.UserID, item.Name, item.DefaultPrice, item.CreatedAt, item.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("inserting service item: %w", err)
	}
	return nil
}

// Update modifies an existing service item.
func (r *Repository) Update(item *models.ServiceItem) error {
	item.UpdatedAt = time.Now()

	result, err := r.db.Exec(`
		UPDATE service_items SET name=$1, default_price=$2, updated_at=$3
		WHERE id=$4 AND user_id=$5`,
		item.Name, item.DefaultPrice, item.UpdatedAt, item.ID, item.UserID,
	)
	if err != nil {
		return fmt.Errorf("updating service item: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("service item not found")
	}
	return nil
}

// Delete removes a service item.
func (r *Repository) Delete(userID, itemID uuid.UUID) error {
	result, err := r.db.Exec(`DELETE FROM service_items WHERE id=$1 AND user_id=$2`, itemID, userID)
	if err != nil {
		return fmt.Errorf("deleting service item: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("service item not found")
	}
	return nil
}

// GetByID returns a single service item scoped to the user.
func (r *Repository) GetByID(userID, itemID uuid.UUID) (*models.ServiceItem, error) {
	item := &models.ServiceItem{}
	err := r.db.QueryRow(`
		SELECT id, user_id, name, default_price, created_at, updated_at
		FROM service_items WHERE id=$1 AND user_id=$2`, itemID, userID,
	).Scan(&item.ID, &item.UserID, &item.Name, &item.DefaultPrice, &item.CreatedAt, &item.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("querying service item: %w", err)
	}
	return item, nil
}
