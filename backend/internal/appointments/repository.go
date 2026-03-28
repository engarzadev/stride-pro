// Package appointments manages appointment scheduling data access and business logic.
package appointments

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/stride-pro/backend/internal/database"
	"github.com/stride-pro/backend/internal/models"
)

// Repository handles appointment persistence.
type Repository struct {
	db *database.DB
}

// NewRepository creates an appointment repository.
func NewRepository(db *database.DB) *Repository {
	return &Repository{db: db}
}

const apptSelectCols = `a.id, a.user_id, a.client_id, a.horse_id, a.barn_id, a.scheduled_at, a.duration, a.status, a.type, a.notes, a.created_at, a.updated_at, c.id, c.first_name, c.last_name, h.id, h.name, b.id, b.name`
const apptJoins = `FROM appointments a LEFT JOIN clients c ON a.client_id = c.id LEFT JOIN horses h ON a.horse_id = h.id LEFT JOIN barns b ON a.barn_id = b.id`

func scanAppointment(scanner interface{ Scan(...interface{}) error }) (*models.Appointment, error) {
	a := &models.Appointment{}
	var clientID uuid.NullUUID
	var clientFirstName, clientLastName sql.NullString
	var horseID uuid.NullUUID
	var horseName sql.NullString
	var barnID uuid.NullUUID
	var barnName sql.NullString
	err := scanner.Scan(
		&a.ID, &a.UserID, &a.ClientID, &a.HorseID, &a.BarnID,
		&a.ScheduledAt, &a.Duration, &a.Status, &a.Type, &a.Notes,
		&a.CreatedAt, &a.UpdatedAt,
		&clientID, &clientFirstName, &clientLastName,
		&horseID, &horseName,
		&barnID, &barnName,
	)
	if err != nil {
		return nil, err
	}
	if clientID.Valid {
		a.Client = &models.Client{
			ID:        clientID.UUID,
			FirstName: clientFirstName.String,
			LastName:  clientLastName.String,
		}
	}
	if horseID.Valid {
		a.Horse = &models.Horse{
			ID:   horseID.UUID,
			Name: horseName.String,
		}
	}
	if barnID.Valid {
		a.Barn = &models.Barn{
			ID:   barnID.UUID,
			Name: barnName.String,
		}
	}
	return a, nil
}

// Create inserts a new appointment.
func (r *Repository) Create(a *models.Appointment) error {
	a.ID = uuid.New()
	a.CreatedAt = time.Now()
	a.UpdatedAt = time.Now()

	_, err := r.db.Exec(`
		INSERT INTO appointments (id, user_id, client_id, horse_id, barn_id, scheduled_at, duration, status, type, notes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
		a.ID, a.UserID, a.ClientID, a.HorseID, a.BarnID,
		a.ScheduledAt, a.Duration, a.Status, a.Type, a.Notes,
		a.CreatedAt, a.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("inserting appointment: %w", err)
	}
	return nil
}

// GetByID returns a single appointment by ID, scoped to the user.
func (r *Repository) GetByID(userID, apptID uuid.UUID) (*models.Appointment, error) {
	a, err := scanAppointment(r.db.QueryRow(
		`SELECT `+apptSelectCols+` `+apptJoins+` WHERE a.id = $1 AND a.user_id = $2`,
		apptID, userID,
	))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("querying appointment: %w", err)
	}
	return a, nil
}

// GetAllByUserID returns all appointments belonging to a user.
func (r *Repository) GetAllByUserID(userID uuid.UUID) ([]models.Appointment, error) {
	rows, err := r.db.Query(
		`SELECT `+apptSelectCols+` `+apptJoins+` WHERE a.user_id = $1 ORDER BY a.scheduled_at DESC`, userID,
	)
	if err != nil {
		return nil, fmt.Errorf("querying appointments: %w", err)
	}
	defer rows.Close()

	var appts []models.Appointment
	for rows.Next() {
		a, err := scanAppointment(rows)
		if err != nil {
			return nil, fmt.Errorf("scanning appointment: %w", err)
		}
		appts = append(appts, *a)
	}
	return appts, rows.Err()
}

// GetByDateRange returns appointments within a date range for the user.
func (r *Repository) GetByDateRange(userID uuid.UUID, start, end time.Time) ([]models.Appointment, error) {
	rows, err := r.db.Query(
		`SELECT `+apptSelectCols+` `+apptJoins+`
		 WHERE a.user_id = $1 AND a.scheduled_at >= $2 AND a.scheduled_at <= $3
		 ORDER BY a.scheduled_at`,
		userID, start, end,
	)
	if err != nil {
		return nil, fmt.Errorf("querying appointments by date range: %w", err)
	}
	defer rows.Close()

	var appts []models.Appointment
	for rows.Next() {
		a, err := scanAppointment(rows)
		if err != nil {
			return nil, fmt.Errorf("scanning appointment: %w", err)
		}
		appts = append(appts, *a)
	}
	return appts, rows.Err()
}

// Update modifies an existing appointment.
func (r *Repository) Update(a *models.Appointment) error {
	a.UpdatedAt = time.Now()
	result, err := r.db.Exec(`
		UPDATE appointments SET client_id=$1, horse_id=$2, barn_id=$3, scheduled_at=$4,
		duration=$5, status=$6, type=$7, notes=$8, updated_at=$9
		WHERE id=$10 AND user_id=$11`,
		a.ClientID, a.HorseID, a.BarnID, a.ScheduledAt,
		a.Duration, a.Status, a.Type, a.Notes, a.UpdatedAt,
		a.ID, a.UserID,
	)
	if err != nil {
		return fmt.Errorf("updating appointment: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("appointment not found")
	}
	return nil
}

// Delete removes an appointment by ID, scoped to the user.
func (r *Repository) Delete(userID, apptID uuid.UUID) error {
	result, err := r.db.Exec("DELETE FROM appointments WHERE id = $1 AND user_id = $2", apptID, userID)
	if err != nil {
		return fmt.Errorf("deleting appointment: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("appointment not found")
	}
	return nil
}
