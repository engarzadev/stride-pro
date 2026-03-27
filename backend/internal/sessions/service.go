package sessions

import (
	"encoding/json"

	"github.com/google/uuid"

	"github.com/stride-pro/backend/internal/models"
	"github.com/stride-pro/backend/pkg/validator"
)

// Service contains business logic for session management.
type Service struct {
	repo *Repository
}

// NewService creates a session service.
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// CreateInput holds data for creating or updating a session.
type CreateInput struct {
	AppointmentID   uuid.UUID `json:"appointment_id"`
	Type            string    `json:"type"`
	BodyZones       []string  `json:"body_zones"`
	Notes           string    `json:"notes"`
	Findings        string    `json:"findings"`
	Recommendations string    `json:"recommendations"`
}

// Validate checks the input for errors.
func (i *CreateInput) Validate() validator.Errors {
	errs := validator.Errors{}
	validator.OneOf(errs, "type", i.Type, []string{"skeletal", "muscular", "soft_tissue", "other"})
	return errs
}

// Create validates and creates a new session.
func (s *Service) Create(userID uuid.UUID, input CreateInput) (*models.Session, error) {
	bodyZonesJSON, _ := json.Marshal(input.BodyZones)

	sess := &models.Session{
		UserID:          userID,
		AppointmentID:   input.AppointmentID,
		Type:            input.Type,
		BodyZones:       bodyZonesJSON,
		Notes:           input.Notes,
		Findings:        input.Findings,
		Recommendations: input.Recommendations,
	}
	if err := s.repo.Create(sess); err != nil {
		return nil, err
	}
	return sess, nil
}

// GetByID returns a session by ID for the given user.
func (s *Service) GetByID(userID, sessionID uuid.UUID) (*models.Session, error) {
	return s.repo.GetByID(userID, sessionID)
}

// GetAll returns all sessions for the given user.
func (s *Service) GetAll(userID uuid.UUID) ([]models.Session, error) {
	return s.repo.GetAllByUserID(userID)
}

// GetByAppointmentID returns all sessions for a specific appointment.
func (s *Service) GetByAppointmentID(userID, appointmentID uuid.UUID) ([]models.Session, error) {
	return s.repo.GetByAppointmentID(userID, appointmentID)
}

// Update modifies an existing session.
func (s *Service) Update(userID, sessionID uuid.UUID, input CreateInput) (*models.Session, error) {
	bodyZonesJSON, _ := json.Marshal(input.BodyZones)

	sess := &models.Session{
		ID:              sessionID,
		UserID:          userID,
		AppointmentID:   input.AppointmentID,
		Type:            input.Type,
		BodyZones:       bodyZonesJSON,
		Notes:           input.Notes,
		Findings:        input.Findings,
		Recommendations: input.Recommendations,
	}
	if err := s.repo.Update(sess); err != nil {
		return nil, err
	}
	return sess, nil
}

// Delete removes a session.
func (s *Service) Delete(userID, sessionID uuid.UUID) error {
	return s.repo.Delete(userID, sessionID)
}
