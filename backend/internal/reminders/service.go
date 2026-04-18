package reminders

import (
	"github.com/google/uuid"

	"github.com/stride-pro/backend/internal/models"
	"github.com/stride-pro/backend/pkg/validator"
)

// Service contains business logic for reminder management.
type Service struct {
	repo *Repository
}

// NewService creates a reminder service.
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// CreateInput holds data for creating a reminder.
type CreateInput struct {
	Title      string `json:"title"`
	DueDate    string `json:"due_date"`
	Category   string `json:"category"`
	Source     string `json:"source"`
	IsComplete bool   `json:"is_complete"`
}

// Validate checks the input for errors.
func (i *CreateInput) Validate() validator.Errors {
	errs := validator.Errors{}
	validator.Required(errs, "title", i.Title)
	validator.Required(errs, "due_date", i.DueDate)
	if i.Source == "" {
		i.Source = "manual"
	}
	return errs
}

// PatchInput holds fields that can be updated on an existing reminder.
type PatchInput struct {
	IsComplete *bool `json:"is_complete"`
}

// GetByHorseID returns all reminders for a horse.
func (s *Service) GetByHorseID(userID, horseID uuid.UUID) ([]models.Reminder, error) {
	list, err := s.repo.GetByHorseID(userID, horseID)
	if err != nil {
		return nil, err
	}
	if list == nil {
		list = []models.Reminder{}
	}
	return list, nil
}

// Create adds a new reminder.
func (s *Service) Create(userID, horseID uuid.UUID, input CreateInput) (*models.Reminder, error) {
	rm := &models.Reminder{
		UserID:     userID,
		HorseID:    horseID,
		Title:      input.Title,
		DueDate:    input.DueDate,
		Category:   input.Category,
		Source:     input.Source,
		IsComplete: input.IsComplete,
	}
	if err := s.repo.Create(rm); err != nil {
		return nil, err
	}
	return rm, nil
}

// Update replaces the editable fields of a reminder.
func (s *Service) Update(userID, id uuid.UUID, input CreateInput) (*models.Reminder, error) {
	rm, err := s.repo.Update(userID, id, input.Title, input.DueDate, input.Category)
	if err != nil {
		return nil, err
	}
	return rm, nil
}

// Patch applies a partial update to a reminder.
func (s *Service) Patch(userID, id uuid.UUID, input PatchInput) (*models.Reminder, error) {
	if input.IsComplete == nil {
		return s.repo.GetByID(userID, id)
	}
	rm, err := s.repo.SetComplete(userID, id, *input.IsComplete)
	if err != nil {
		return nil, err
	}
	return rm, nil
}

// Delete removes a reminder.
func (s *Service) Delete(userID, id uuid.UUID) error {
	return s.repo.Delete(userID, id)
}
