// Package sessions manages care session data access and business logic.
package sessions

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/stride-pro/backend/internal/database"
	"github.com/stride-pro/backend/internal/models"
)

// Repository handles session persistence.
type Repository struct {
	db *database.DB
}

// NewRepository creates a session repository.
func NewRepository(db *database.DB) *Repository {
	return &Repository{db: db}
}

const sessionColumns = `id, user_id, appointment_id, type, body_zones, notes, findings, recommendations, created_at, updated_at`

func scanSession(scanner interface{ Scan(...interface{}) error }) (*models.Session, error) {
	s := &models.Session{}
	err := scanner.Scan(
		&s.ID, &s.UserID, &s.AppointmentID, &s.Type, &s.BodyZones,
		&s.Notes, &s.Findings, &s.Recommendations, &s.CreatedAt, &s.UpdatedAt,
	)
	return s, err
}

// Create inserts a new session.
func (r *Repository) Create(s *models.Session) error {
	s.ID = uuid.New()
	s.CreatedAt = time.Now()
	s.UpdatedAt = time.Now()

	_, err := r.db.Exec(`
		INSERT INTO sessions (`+sessionColumns+`)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		s.ID, s.UserID, s.AppointmentID, s.Type, s.BodyZones,
		s.Notes, s.Findings, s.Recommendations, s.CreatedAt, s.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("inserting session: %w", err)
	}
	return nil
}

// GetByID returns a single session by ID, scoped to the user.
func (r *Repository) GetByID(userID, sessionID uuid.UUID) (*models.Session, error) {
	s, err := scanSession(r.db.QueryRow(
		`SELECT `+sessionColumns+` FROM sessions WHERE id = $1 AND user_id = $2`,
		sessionID, userID,
	))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("querying session: %w", err)
	}
	return s, nil
}

// GetAllByUserID returns all sessions belonging to a user.
func (r *Repository) GetAllByUserID(userID uuid.UUID) ([]models.Session, error) {
	rows, err := r.db.Query(
		`SELECT `+sessionColumns+` FROM sessions WHERE user_id = $1 ORDER BY created_at DESC`, userID,
	)
	if err != nil {
		return nil, fmt.Errorf("querying sessions: %w", err)
	}
	defer rows.Close()

	var sessions []models.Session
	for rows.Next() {
		s, err := scanSession(rows)
		if err != nil {
			return nil, fmt.Errorf("scanning session: %w", err)
		}
		sessions = append(sessions, *s)
	}
	return sessions, rows.Err()
}

// GetByAppointmentID returns all sessions for a specific appointment.
func (r *Repository) GetByAppointmentID(userID, appointmentID uuid.UUID) ([]models.Session, error) {
	rows, err := r.db.Query(
		`SELECT `+sessionColumns+` FROM sessions WHERE user_id = $1 AND appointment_id = $2 ORDER BY created_at`,
		userID, appointmentID,
	)
	if err != nil {
		return nil, fmt.Errorf("querying sessions by appointment: %w", err)
	}
	defer rows.Close()

	var sessions []models.Session
	for rows.Next() {
		s, err := scanSession(rows)
		if err != nil {
			return nil, fmt.Errorf("scanning session: %w", err)
		}
		sessions = append(sessions, *s)
	}
	return sessions, rows.Err()
}

// Update modifies an existing session.
func (r *Repository) Update(s *models.Session) error {
	s.UpdatedAt = time.Now()
	result, err := r.db.Exec(`
		UPDATE sessions SET appointment_id=$1, type=$2, body_zones=$3, notes=$4,
		findings=$5, recommendations=$6, updated_at=$7
		WHERE id=$8 AND user_id=$9`,
		s.AppointmentID, s.Type, s.BodyZones, s.Notes,
		s.Findings, s.Recommendations, s.UpdatedAt,
		s.ID, s.UserID,
	)
	if err != nil {
		return fmt.Errorf("updating session: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("session not found")
	}
	return nil
}

// Delete removes a session by ID, scoped to the user.
func (r *Repository) Delete(userID, sessionID uuid.UUID) error {
	result, err := r.db.Exec("DELETE FROM sessions WHERE id = $1 AND user_id = $2", sessionID, userID)
	if err != nil {
		return fmt.Errorf("deleting session: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("session not found")
	}
	return nil
}
