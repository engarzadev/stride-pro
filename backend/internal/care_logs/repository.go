// Package care_logs manages care log data access and business logic.
package care_logs

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/stride-pro/backend/internal/database"
	"github.com/stride-pro/backend/internal/models"
)

// Repository handles care log persistence.
type Repository struct {
	db *database.DB
}

// NewRepository creates a care log repository.
func NewRepository(db *database.DB) *Repository {
	return &Repository{db: db}
}

// GetByHorseID returns all care logs for a horse, newest first.
func (r *Repository) GetByHorseID(userID, horseID uuid.UUID) ([]models.CareLog, error) {
	rows, err := r.db.Query(`
		SELECT id, user_id, horse_id, date, category, notes, created_at, updated_at
		FROM care_logs
		WHERE user_id = $1 AND horse_id = $2
		ORDER BY date DESC, created_at DESC`,
		userID, horseID,
	)
	if err != nil {
		return nil, fmt.Errorf("querying care logs: %w", err)
	}
	defer rows.Close()

	var logs []models.CareLog
	for rows.Next() {
		var cl models.CareLog
		if err := rows.Scan(&cl.ID, &cl.UserID, &cl.HorseID, &cl.Date, &cl.Category, &cl.Notes, &cl.CreatedAt, &cl.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scanning care log: %w", err)
		}
		logs = append(logs, cl)
	}
	return logs, rows.Err()
}

// GetByID returns a single care log, scoped to the user.
func (r *Repository) GetByID(userID, id uuid.UUID) (*models.CareLog, error) {
	cl := &models.CareLog{}
	err := r.db.QueryRow(`
		SELECT id, user_id, horse_id, date, category, notes, created_at, updated_at
		FROM care_logs WHERE id = $1 AND user_id = $2`,
		id, userID,
	).Scan(&cl.ID, &cl.UserID, &cl.HorseID, &cl.Date, &cl.Category, &cl.Notes, &cl.CreatedAt, &cl.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("querying care log: %w", err)
	}
	return cl, nil
}

// Create inserts a new care log.
func (r *Repository) Create(cl *models.CareLog) error {
	cl.ID = uuid.New()
	cl.CreatedAt = time.Now()
	cl.UpdatedAt = time.Now()

	_, err := r.db.Exec(`
		INSERT INTO care_logs (id, user_id, horse_id, date, category, notes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		cl.ID, cl.UserID, cl.HorseID, cl.Date, cl.Category, cl.Notes, cl.CreatedAt, cl.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("inserting care log: %w", err)
	}
	return nil
}

// Update modifies an existing care log.
func (r *Repository) Update(cl *models.CareLog) error {
	cl.UpdatedAt = time.Now()
	result, err := r.db.Exec(`
		UPDATE care_logs SET date=$1, category=$2, notes=$3, updated_at=$4
		WHERE id=$5 AND user_id=$6`,
		cl.Date, cl.Category, cl.Notes, cl.UpdatedAt, cl.ID, cl.UserID,
	)
	if err != nil {
		return fmt.Errorf("updating care log: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("care log not found")
	}
	return nil
}

// Delete removes a care log by ID, scoped to the user.
func (r *Repository) Delete(userID, id uuid.UUID) error {
	result, err := r.db.Exec("DELETE FROM care_logs WHERE id = $1 AND user_id = $2", id, userID)
	if err != nil {
		return fmt.Errorf("deleting care log: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("care log not found")
	}
	return nil
}
