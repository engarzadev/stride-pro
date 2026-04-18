package reminders

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/stride-pro/backend/internal/database"
	"github.com/stride-pro/backend/internal/models"
)

// Repository handles reminder persistence.
type Repository struct {
	db *database.DB
}

// NewRepository creates a reminder repository.
func NewRepository(db *database.DB) *Repository {
	return &Repository{db: db}
}

const scanCols = `id, user_id, horse_id, title, due_date, category, source, is_complete, created_at, updated_at`

func scanReminder(row interface{ Scan(...any) error }) (*models.Reminder, error) {
	rm := &models.Reminder{}
	err := row.Scan(&rm.ID, &rm.UserID, &rm.HorseID, &rm.Title, &rm.DueDate, &rm.Category, &rm.Source, &rm.IsComplete, &rm.CreatedAt, &rm.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("scanning reminder: %w", err)
	}
	return rm, nil
}

// GetByHorseID returns all reminders for a horse, active first then completed, each sorted by due date.
func (r *Repository) GetByHorseID(userID, horseID uuid.UUID) ([]models.Reminder, error) {
	rows, err := r.db.Query(
		`SELECT `+scanCols+` FROM reminders
		 WHERE user_id = $1 AND horse_id = $2
		 ORDER BY is_complete ASC, due_date ASC`,
		userID, horseID,
	)
	if err != nil {
		return nil, fmt.Errorf("querying reminders: %w", err)
	}
	defer rows.Close()

	var list []models.Reminder
	for rows.Next() {
		rm, err := scanReminder(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, *rm)
	}
	return list, rows.Err()
}

// GetByID returns a single reminder scoped to the user.
func (r *Repository) GetByID(userID, id uuid.UUID) (*models.Reminder, error) {
	row := r.db.QueryRow(
		`SELECT `+scanCols+` FROM reminders WHERE id = $1 AND user_id = $2`,
		id, userID,
	)
	return scanReminder(row)
}

// Create inserts a new reminder.
func (r *Repository) Create(rm *models.Reminder) error {
	rm.ID = uuid.New()
	rm.CreatedAt = time.Now()
	rm.UpdatedAt = time.Now()

	_, err := r.db.Exec(
		`INSERT INTO reminders (id, user_id, horse_id, title, due_date, category, source, is_complete, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		rm.ID, rm.UserID, rm.HorseID, rm.Title, rm.DueDate, rm.Category, rm.Source, rm.IsComplete, rm.CreatedAt, rm.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("inserting reminder: %w", err)
	}
	return nil
}

// Update replaces the mutable fields of a reminder.
func (r *Repository) Update(userID, id uuid.UUID, title, dueDate, category string) (*models.Reminder, error) {
	_, err := r.db.Exec(
		`UPDATE reminders SET title=$1, due_date=$2, category=$3, updated_at=$4 WHERE id=$5 AND user_id=$6`,
		title, dueDate, category, time.Now(), id, userID,
	)
	if err != nil {
		return nil, fmt.Errorf("updating reminder: %w", err)
	}
	return r.GetByID(userID, id)
}

// SetComplete toggles the is_complete flag on a reminder.
func (r *Repository) SetComplete(userID, id uuid.UUID, complete bool) (*models.Reminder, error) {
	_, err := r.db.Exec(
		`UPDATE reminders SET is_complete=$1, updated_at=$2 WHERE id=$3 AND user_id=$4`,
		complete, time.Now(), id, userID,
	)
	if err != nil {
		return nil, fmt.Errorf("updating reminder: %w", err)
	}
	return r.GetByID(userID, id)
}

// Delete removes a reminder by ID scoped to the user.
func (r *Repository) Delete(userID, id uuid.UUID) error {
	result, err := r.db.Exec("DELETE FROM reminders WHERE id = $1 AND user_id = $2", id, userID)
	if err != nil {
		return fmt.Errorf("deleting reminder: %w", err)
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return fmt.Errorf("reminder not found")
	}
	return nil
}
