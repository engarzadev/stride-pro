package appointments

import (
	"fmt"
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
	ClientID uuid.UUID  `json:"client_id"`
	HorseID  uuid.UUID  `json:"horse_id"`
	BarnID   *uuid.UUID `json:"barn_id"`
	Date     string     `json:"date"`
	Time     string     `json:"time"`
	Duration int        `json:"duration"`
	Status   string     `json:"status"`
	Type     string     `json:"type"`
	Notes    string     `json:"notes"`
}

// scheduledAt combines the Date and Time strings into a time.Time.
func (i *CreateInput) scheduledAt() (time.Time, error) {
	t := i.Time
	if t == "" {
		t = "00:00"
	}
	return time.ParseInLocation("2006-01-02 15:04", i.Date+" "+t, time.UTC)
}

// Validate checks the input for errors.
func (i *CreateInput) Validate() validator.Errors {
	errs := validator.Errors{}
	validator.Required(errs, "date", i.Date)
	validator.MinValue(errs, "duration", i.Duration, 1)
	validator.OneOf(errs, "status", i.Status, []string{"scheduled", "completed", "cancelled", "no-show"})
	validator.Required(errs, "type", i.Type)
	return errs
}

// populateDateFields sets the Date and Time string fields from ScheduledAt.
func populateDateFields(a *models.Appointment) {
	a.Date = a.ScheduledAt.UTC().Format("2006-01-02")
	a.Time = a.ScheduledAt.UTC().Format("15:04")
}

// Create validates and creates a new appointment.
func (s *Service) Create(userID uuid.UUID, input CreateInput) (*models.Appointment, error) {
	scheduledAt, err := input.scheduledAt()
	if err != nil {
		return nil, fmt.Errorf("parsing date/time: %w", err)
	}
	a := &models.Appointment{
		UserID:      userID,
		ClientID:    input.ClientID,
		HorseID:     input.HorseID,
		BarnID:      models.PtrToNullUUID(input.BarnID),
		ScheduledAt: scheduledAt,
		Duration:    input.Duration,
		Status:      input.Status,
		Type:        input.Type,
		Notes:       input.Notes,
	}
	if err := s.repo.Create(a); err != nil {
		return nil, err
	}
	populateDateFields(a)
	return a, nil
}

// GetByID returns an appointment by ID for the given user.
func (s *Service) GetByID(userID, apptID uuid.UUID) (*models.Appointment, error) {
	a, err := s.repo.GetByID(userID, apptID)
	if err != nil || a == nil {
		return a, err
	}
	populateDateFields(a)
	return a, nil
}

// GetAll returns all appointments for the given user.
func (s *Service) GetAll(userID uuid.UUID) ([]models.Appointment, error) {
	appts, err := s.repo.GetAllByUserID(userID)
	if err != nil {
		return nil, err
	}
	for i := range appts {
		populateDateFields(&appts[i])
	}
	return appts, nil
}

// GetByDateRange returns appointments within a date range.
func (s *Service) GetByDateRange(userID uuid.UUID, start, end time.Time) ([]models.Appointment, error) {
	appts, err := s.repo.GetByDateRange(userID, start, end)
	if err != nil {
		return nil, err
	}
	for i := range appts {
		populateDateFields(&appts[i])
	}
	return appts, nil
}

// Update modifies an existing appointment.
func (s *Service) Update(userID, apptID uuid.UUID, input CreateInput) (*models.Appointment, error) {
	scheduledAt, err := input.scheduledAt()
	if err != nil {
		return nil, fmt.Errorf("parsing date/time: %w", err)
	}
	a := &models.Appointment{
		ID:          apptID,
		UserID:      userID,
		ClientID:    input.ClientID,
		HorseID:     input.HorseID,
		BarnID:      models.PtrToNullUUID(input.BarnID),
		ScheduledAt: scheduledAt,
		Duration:    input.Duration,
		Status:      input.Status,
		Type:        input.Type,
		Notes:       input.Notes,
	}
	if err := s.repo.Update(a); err != nil {
		return nil, err
	}
	populateDateFields(a)
	return a, nil
}

// Delete removes an appointment.
func (s *Service) Delete(userID, apptID uuid.UUID) error {
	return s.repo.Delete(userID, apptID)
}
