package appointments

import (
	"time"

	"github.com/google/uuid"

	"github.com/stride-pro/backend/internal/models"
	"github.com/stride-pro/backend/pkg/validator"
)

// Service contains business logic for appointment management.
type Service struct {
	repo *Repository
}

// NewService creates an appointment service.
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// CreateInput holds data for creating or updating an appointment.
type CreateInput struct {
	ClientID    uuid.UUID  `json:"client_id"`
	HorseID     uuid.UUID  `json:"horse_id"`
	BarnID      *uuid.UUID `json:"barn_id"`
	ScheduledAt time.Time  `json:"scheduled_at"`
	Duration    int        `json:"duration"`
	Status      string     `json:"status"`
	Type        string     `json:"type"`
	Notes       string     `json:"notes"`
}

// Validate checks the input for errors.
func (i *CreateInput) Validate() validator.Errors {
	errs := validator.Errors{}
	if i.ScheduledAt.IsZero() {
		errs["scheduled_at"] = "scheduled_at is required"
	}
	validator.MinValue(errs, "duration", i.Duration, 1)
	validator.OneOf(errs, "status", i.Status, []string{"scheduled", "completed", "cancelled", "no-show"})
	validator.Required(errs, "type", i.Type)
	return errs
}

// Create validates and creates a new appointment.
func (s *Service) Create(userID uuid.UUID, input CreateInput) (*models.Appointment, error) {
	a := &models.Appointment{
		UserID:      userID,
		ClientID:    input.ClientID,
		HorseID:     input.HorseID,
		BarnID:      models.PtrToNullUUID(input.BarnID),
		ScheduledAt: input.ScheduledAt,
		Duration:    input.Duration,
		Status:      input.Status,
		Type:        input.Type,
		Notes:       input.Notes,
	}
	if err := s.repo.Create(a); err != nil {
		return nil, err
	}
	return a, nil
}

// GetByID returns an appointment by ID for the given user.
func (s *Service) GetByID(userID, apptID uuid.UUID) (*models.Appointment, error) {
	return s.repo.GetByID(userID, apptID)
}

// GetAll returns all appointments for the given user.
func (s *Service) GetAll(userID uuid.UUID) ([]models.Appointment, error) {
	return s.repo.GetAllByUserID(userID)
}

// GetByDateRange returns appointments within a date range.
func (s *Service) GetByDateRange(userID uuid.UUID, start, end time.Time) ([]models.Appointment, error) {
	return s.repo.GetByDateRange(userID, start, end)
}

// Update modifies an existing appointment.
func (s *Service) Update(userID, apptID uuid.UUID, input CreateInput) (*models.Appointment, error) {
	a := &models.Appointment{
		ID:          apptID,
		UserID:      userID,
		ClientID:    input.ClientID,
		HorseID:     input.HorseID,
		BarnID:      models.PtrToNullUUID(input.BarnID),
		ScheduledAt: input.ScheduledAt,
		Duration:    input.Duration,
		Status:      input.Status,
		Type:        input.Type,
		Notes:       input.Notes,
	}
	if err := s.repo.Update(a); err != nil {
		return nil, err
	}
	return a, nil
}

// Delete removes an appointment.
func (s *Service) Delete(userID, apptID uuid.UUID) error {
	return s.repo.Delete(userID, apptID)
}
