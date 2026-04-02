package appointments

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/google/uuid"

	"github.com/stride-pro/backend/internal/auth"
	"github.com/stride-pro/backend/internal/models"
	"github.com/stride-pro/backend/internal/notifications"
	"github.com/stride-pro/backend/pkg/validator"
)

// ErrConflict is returned when a new appointment overlaps an existing one.
var ErrConflict = errors.New("appointment conflicts with an existing appointment")

// Service contains business logic for appointment management.
type Service struct {
	repo          *Repository
	notifications *notifications.Service
	authService   *auth.Service
}

// NewService creates an appointment service.
func NewService(repo *Repository, notifSvc *notifications.Service, authSvc *auth.Service) *Service {
	return &Service{repo: repo, notifications: notifSvc, authService: authSvc}
}

// CreateInput holds data for creating or updating an appointment.
type CreateInput struct {
	ClientID   uuid.UUID  `json:"client_id"`
	HorseID    uuid.UUID  `json:"horse_id"`
	BarnID     *uuid.UUID `json:"barn_id"`
	Date       string     `json:"date"`
	Time       string     `json:"time"`
	Duration   int        `json:"duration"`
	TravelTime int        `json:"travel_time"`
	Status     string     `json:"status"`
	Type       string     `json:"type"`
	Notes      string     `json:"notes"`
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

	conflict, err := s.repo.HasConflict(userID, uuid.Nil, scheduledAt, input.Duration, input.TravelTime)
	if err != nil {
		return nil, fmt.Errorf("checking conflicts: %w", err)
	}
	if conflict {
		return nil, ErrConflict
	}

	a := &models.Appointment{
		UserID:      userID,
		ClientID:    input.ClientID,
		HorseID:     input.HorseID,
		BarnID:      models.PtrToNullUUID(input.BarnID),
		ScheduledAt: scheduledAt,
		Duration:    input.Duration,
		TravelTime:  input.TravelTime,
		Status:      input.Status,
		Type:        input.Type,
		Notes:       input.Notes,
	}
	if err := s.repo.Create(a); err != nil {
		return nil, err
	}
	populateDateFields(a)
	s.sendBookingConfirmation(userID, a)
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

	conflict, err := s.repo.HasConflict(userID, apptID, scheduledAt, input.Duration, input.TravelTime)
	if err != nil {
		return nil, fmt.Errorf("checking conflicts: %w", err)
	}
	if conflict {
		return nil, ErrConflict
	}

	a := &models.Appointment{
		ID:          apptID,
		UserID:      userID,
		ClientID:    input.ClientID,
		HorseID:     input.HorseID,
		BarnID:      models.PtrToNullUUID(input.BarnID),
		ScheduledAt: scheduledAt,
		Duration:    input.Duration,
		TravelTime:  input.TravelTime,
		Status:      input.Status,
		Type:        input.Type,
		Notes:       input.Notes,
	}
	if err := s.repo.Update(a); err != nil {
		return nil, err
	}
	populateDateFields(a)
	s.sendBookingConfirmation(userID, a)
	return a, nil
}

// Delete removes an appointment.
func (s *Service) Delete(userID, apptID uuid.UUID) error {
	return s.repo.Delete(userID, apptID)
}

// sendBookingConfirmation sends a booking confirmation notification to the client.
// It loads the full appointment to get client/horse data and is non-fatal on error.
func (s *Service) sendBookingConfirmation(userID uuid.UUID, a *models.Appointment) {
	if s.notifications == nil {
		return
	}
	full, err := s.repo.GetByID(userID, a.ID)
	if err != nil || full == nil || full.Client == nil || full.Client.Email == "" {
		return
	}
	populateDateFields(full)

	horseName := ""
	if full.Horse != nil {
		horseName = full.Horse.Name
	}

	providerName := ""
	if s.authService != nil {
		if user, err := s.authService.GetUserByID(userID); err == nil {
			providerName = user.FirstName + " " + user.LastName
		}
	}

	subject, body, err := notifications.RenderTemplate("booking_confirmation", map[string]string{
		"client_name":      full.Client.FirstName + " " + full.Client.LastName,
		"date":             full.Date,
		"time":             full.Time,
		"horse_name":       horseName,
		"appointment_type": full.Type,
		"duration":         strconv.Itoa(full.Duration),
		"provider_name":    providerName,
	})
	if err != nil {
		log.Printf("rendering booking confirmation template: %v", err)
		return
	}

	n := &models.Notification{
		UserID:    userID,
		Type:      "booking_confirmation",
		Channel:   "email",
		Recipient: full.Client.Email,
		Subject:   subject,
		Body:      body,
	}
	if err := s.notifications.Send(n); err != nil {
		log.Printf("sending booking confirmation for appointment %s: %v", a.ID, err)
	}
}
