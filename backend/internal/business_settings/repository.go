// Package business_settings manages user business profile data.
package business_settings

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/stride-pro/backend/internal/database"
	"github.com/stride-pro/backend/internal/models"
)

// Repository handles business settings persistence.
type Repository struct {
	db *database.DB
}

// NewRepository creates a business settings repository.
func NewRepository(db *database.DB) *Repository {
	return &Repository{db: db}
}

// Get returns the business settings for a user, or nil if not set.
func (r *Repository) Get(userID uuid.UUID) (*models.BusinessSettings, error) {
	bs := &models.BusinessSettings{}
	err := r.db.QueryRow(`
		SELECT id, user_id, business_name, email, phone, address, invoice_message, created_at, updated_at
		FROM business_settings WHERE user_id = $1`, userID,
	).Scan(&bs.ID, &bs.UserID, &bs.BusinessName, &bs.Email, &bs.Phone, &bs.Address, &bs.InvoiceMessage, &bs.CreatedAt, &bs.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("querying business settings: %w", err)
	}
	return bs, nil
}

// Upsert creates or updates the business settings for a user.
func (r *Repository) Upsert(bs *models.BusinessSettings) error {
	now := time.Now()
	bs.UpdatedAt = now

	err := r.db.QueryRow(`
		INSERT INTO business_settings (id, user_id, business_name, email, phone, address, invoice_message, created_at, updated_at)
		VALUES (uuid_generate_v4(), $1, $2, $3, $4, $5, $6, NOW(), $7)
		ON CONFLICT (user_id) DO UPDATE SET
			business_name = EXCLUDED.business_name,
			email = EXCLUDED.email,
			phone = EXCLUDED.phone,
			address = EXCLUDED.address,
			invoice_message = EXCLUDED.invoice_message,
			updated_at = EXCLUDED.updated_at
		RETURNING id, created_at`,
		bs.UserID, bs.BusinessName, bs.Email, bs.Phone, bs.Address, bs.InvoiceMessage, now,
	).Scan(&bs.ID, &bs.CreatedAt)
	if err != nil {
		return fmt.Errorf("upserting business settings: %w", err)
	}
	return nil
}
