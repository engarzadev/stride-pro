// Package horses manages horse data access and business logic.
package horses

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/stride-pro/backend/internal/database"
	"github.com/stride-pro/backend/internal/models"
)

// Repository handles horse persistence.
type Repository struct {
	db *database.DB
}

// NewRepository creates a horse repository.
func NewRepository(db *database.DB) *Repository {
	return &Repository{db: db}
}

const horseSelectCols = `h.id, h.user_id, h.client_id, h.barn_id, h.name, h.breed, h.age, h.gender, h.color, h.weight, h.notes, h.created_at, h.updated_at, c.id, c.first_name, c.last_name, b.id, b.name`
const horseJoins = `FROM horses h LEFT JOIN clients c ON h.client_id = c.id LEFT JOIN barns b ON h.barn_id = b.id`

func scanHorse(scanner interface{ Scan(...interface{}) error }) (*models.Horse, error) {
	h := &models.Horse{}
	var clientID uuid.NullUUID
	var clientFirstName, clientLastName sql.NullString
	var barnID uuid.NullUUID
	var barnName sql.NullString
	err := scanner.Scan(
		&h.ID, &h.UserID, &h.ClientID, &h.BarnID, &h.Name, &h.Breed,
		&h.Age, &h.Gender, &h.Color, &h.Weight, &h.Notes, &h.CreatedAt, &h.UpdatedAt,
		&clientID, &clientFirstName, &clientLastName,
		&barnID, &barnName,
	)
	if err != nil {
		return nil, err
	}
	if clientID.Valid {
		h.Client = &models.Client{
			ID:        clientID.UUID,
			FirstName: clientFirstName.String,
			LastName:  clientLastName.String,
		}
	}
	if barnID.Valid {
		h.Barn = &models.Barn{
			ID:   barnID.UUID,
			Name: barnName.String,
		}
	}
	return h, nil
}

// Create inserts a new horse.
func (r *Repository) Create(h *models.Horse) error {
	h.ID = uuid.New()
	h.CreatedAt = time.Now()
	h.UpdatedAt = time.Now()

	_, err := r.db.Exec(`
		INSERT INTO horses (id, user_id, client_id, barn_id, name, breed, age, gender, color, weight, notes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`,
		h.ID, h.UserID, h.ClientID, h.BarnID, h.Name, h.Breed,
		h.Age, h.Gender, h.Color, h.Weight, h.Notes, h.CreatedAt, h.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("inserting horse: %w", err)
	}
	return nil
}

// GetByID returns a single horse by ID, scoped to the user.
func (r *Repository) GetByID(userID, horseID uuid.UUID) (*models.Horse, error) {
	h, err := scanHorse(r.db.QueryRow(
		`SELECT `+horseSelectCols+` `+horseJoins+` WHERE h.id = $1 AND h.user_id = $2`,
		horseID, userID,
	))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("querying horse: %w", err)
	}
	return h, nil
}

// GetAllByUserID returns all horses belonging to a user.
func (r *Repository) GetAllByUserID(userID uuid.UUID) ([]models.Horse, error) {
	rows, err := r.db.Query(
		`SELECT `+horseSelectCols+` `+horseJoins+` WHERE h.user_id = $1 ORDER BY h.name`, userID,
	)
	if err != nil {
		return nil, fmt.Errorf("querying horses: %w", err)
	}
	defer rows.Close()

	var horses []models.Horse
	for rows.Next() {
		h, err := scanHorse(rows)
		if err != nil {
			return nil, fmt.Errorf("scanning horse: %w", err)
		}
		horses = append(horses, *h)
	}
	return horses, rows.Err()
}

// GetByClientID returns all horses for a specific client.
func (r *Repository) GetByClientID(userID, clientID uuid.UUID) ([]models.Horse, error) {
	rows, err := r.db.Query(
		`SELECT `+horseSelectCols+` `+horseJoins+` WHERE h.user_id = $1 AND h.client_id = $2 ORDER BY h.name`,
		userID, clientID,
	)
	if err != nil {
		return nil, fmt.Errorf("querying horses by client: %w", err)
	}
	defer rows.Close()

	var horses []models.Horse
	for rows.Next() {
		h, err := scanHorse(rows)
		if err != nil {
			return nil, fmt.Errorf("scanning horse: %w", err)
		}
		horses = append(horses, *h)
	}
	return horses, rows.Err()
}

// GetByBarnID returns all horses at a specific barn.
func (r *Repository) GetByBarnID(userID, barnID uuid.UUID) ([]models.Horse, error) {
	rows, err := r.db.Query(
		`SELECT `+horseSelectCols+` `+horseJoins+` WHERE h.user_id = $1 AND h.barn_id = $2 ORDER BY h.name`,
		userID, barnID,
	)
	if err != nil {
		return nil, fmt.Errorf("querying horses by barn: %w", err)
	}
	defer rows.Close()

	var horses []models.Horse
	for rows.Next() {
		h, err := scanHorse(rows)
		if err != nil {
			return nil, fmt.Errorf("scanning horse: %w", err)
		}
		horses = append(horses, *h)
	}
	return horses, rows.Err()
}

// Update modifies an existing horse.
func (r *Repository) Update(h *models.Horse) error {
	h.UpdatedAt = time.Now()
	result, err := r.db.Exec(`
		UPDATE horses SET client_id=$1, barn_id=$2, name=$3, breed=$4, age=$5, gender=$6,
		color=$7, weight=$8, notes=$9, updated_at=$10
		WHERE id=$11 AND user_id=$12`,
		h.ClientID, h.BarnID, h.Name, h.Breed, h.Age, h.Gender,
		h.Color, h.Weight, h.Notes, h.UpdatedAt, h.ID, h.UserID,
	)
	if err != nil {
		return fmt.Errorf("updating horse: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("horse not found")
	}
	return nil
}

// CountByUserID returns the number of horses belonging to a user.
func (r *Repository) CountByUserID(userID uuid.UUID) (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM horses WHERE user_id = $1", userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("counting horses: %w", err)
	}
	return count, nil
}

// Delete removes a horse by ID, scoped to the user.
func (r *Repository) Delete(userID, horseID uuid.UUID) error {
	result, err := r.db.Exec("DELETE FROM horses WHERE id = $1 AND user_id = $2", horseID, userID)
	if err != nil {
		return fmt.Errorf("deleting horse: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("horse not found")
	}
	return nil
}
