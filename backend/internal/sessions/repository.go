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

const sessionSelectCols = `s.id, s.user_id, s.appointment_id, s.type, s.body_zones, s.notes, s.findings, s.recommendations, s.created_at, s.updated_at, a.id, c.id, c.first_name, c.last_name, h.id, h.name`
const sessionJoins = `FROM sessions s LEFT JOIN appointments a ON s.appointment_id = a.id LEFT JOIN clients c ON a.client_id = c.id LEFT JOIN horses h ON a.horse_id = h.id`

func scanSession(scanner interface{ Scan(...interface{}) error }) (*models.Session, error) {
	s := &models.Session{}
	var apptID uuid.NullUUID
	var clientID uuid.NullUUID
	var clientFirstName, clientLastName sql.NullString
	var horseID uuid.NullUUID
	var horseName sql.NullString
	err := scanner.Scan(
		&s.ID, &s.UserID, &s.AppointmentID, &s.Type, &s.BodyZones,
		&s.Notes, &s.Findings, &s.Recommendations, &s.CreatedAt, &s.UpdatedAt,
		&apptID, &clientID, &clientFirstName, &clientLastName, &horseID, &horseName,
	)
	if err != nil {
		return nil, err
	}
	if apptID.Valid {
		s.Appointment = &models.Appointment{ID: apptID.UUID}
		if clientID.Valid {
			s.Appointment.Client = &models.Client{
				ID:        clientID.UUID,
				FirstName: clientFirstName.String,
				LastName:  clientLastName.String,
			}
		}
		if horseID.Valid {
			s.Appointment.Horse = &models.Horse{
				ID:   horseID.UUID,
				Name: horseName.String,
			}
		}
	}
	return s, nil
}

// Create inserts a new session.
func (r *Repository) Create(s *models.Session) error {
	s.ID = uuid.New()
	s.CreatedAt = time.Now()
	s.UpdatedAt = time.Now()

	_, err := r.db.Exec(`
		INSERT INTO sessions (id, user_id, appointment_id, type, body_zones, notes, findings, recommendations, created_at, updated_at)
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
		`SELECT `+sessionSelectCols+` `+sessionJoins+` WHERE s.id = $1 AND s.user_id = $2`,
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
		`SELECT `+sessionSelectCols+` `+sessionJoins+` WHERE s.user_id = $1 ORDER BY s.created_at DESC`, userID,
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
		`SELECT `+sessionSelectCols+` `+sessionJoins+` WHERE s.user_id = $1 AND s.appointment_id = $2 ORDER BY s.created_at`,
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
