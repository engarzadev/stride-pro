// Package barns manages barn/location data access and business logic.
package barns

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/stride-pro/backend/internal/database"
	"github.com/stride-pro/backend/internal/models"
)

// Repository handles barn persistence.
type Repository struct {
	db *database.DB
}

// NewRepository creates a barn repository.
func NewRepository(db *database.DB) *Repository {
	return &Repository{db: db}
}

const barnColumns = `id, user_id, name, address, phone, email, notes, created_at, updated_at`

func scanBarn(scanner interface{ Scan(...interface{}) error }) (*models.Barn, error) {
	b := &models.Barn{}
	err := scanner.Scan(&b.ID, &b.UserID, &b.Name, &b.Address, &b.Phone, &b.Email, &b.Notes, &b.CreatedAt, &b.UpdatedAt)
	return b, err
}

// Create inserts a new barn.
func (r *Repository) Create(b *models.Barn) error {
	b.ID = uuid.New()
	b.CreatedAt = time.Now()
	b.UpdatedAt = time.Now()

	_, err := r.db.Exec(`
		INSERT INTO barns (`+barnColumns+`)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		b.ID, b.UserID, b.Name, b.Address, b.Phone, b.Email, b.Notes, b.CreatedAt, b.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("inserting barn: %w", err)
	}
	return nil
}

// GetByID returns a single barn by ID, scoped to the user.
func (r *Repository) GetByID(userID, barnID uuid.UUID) (*models.Barn, error) {
	b, err := scanBarn(r.db.QueryRow(
		`SELECT `+barnColumns+` FROM barns WHERE id = $1 AND user_id = $2`, barnID, userID,
	))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("querying barn: %w", err)
	}
	return b, nil
}

// GetAllByUserID returns all barns belonging to a user.
func (r *Repository) GetAllByUserID(userID uuid.UUID) ([]models.Barn, error) {
	rows, err := r.db.Query(
		`SELECT `+barnColumns+` FROM barns WHERE user_id = $1 ORDER BY name`, userID,
	)
	if err != nil {
		return nil, fmt.Errorf("querying barns: %w", err)
	}
	defer rows.Close()

	var barns []models.Barn
	for rows.Next() {
		b, err := scanBarn(rows)
		if err != nil {
			return nil, fmt.Errorf("scanning barn: %w", err)
		}
		barns = append(barns, *b)
	}
	return barns, rows.Err()
}

// Update modifies an existing barn.
func (r *Repository) Update(b *models.Barn) error {
	b.UpdatedAt = time.Now()
	result, err := r.db.Exec(`
		UPDATE barns SET name=$1, address=$2, phone=$3, email=$4, notes=$5, updated_at=$6
		WHERE id=$7 AND user_id=$8`,
		b.Name, b.Address, b.Phone, b.Email, b.Notes, b.UpdatedAt, b.ID, b.UserID,
	)
	if err != nil {
		return fmt.Errorf("updating barn: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("barn not found")
	}
	return nil
}

// Delete removes a barn by ID, scoped to the user.
func (r *Repository) Delete(userID, barnID uuid.UUID) error {
	result, err := r.db.Exec("DELETE FROM barns WHERE id = $1 AND user_id = $2", barnID, userID)
	if err != nil {
		return fmt.Errorf("deleting barn: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("barn not found")
	}
	return nil
}
