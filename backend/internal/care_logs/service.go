package care_logs

import (
	"github.com/google/uuid"

	"github.com/stride-pro/backend/internal/models"
	"github.com/stride-pro/backend/internal/subscriptions"
	"github.com/stride-pro/backend/pkg/validator"
)

var validCategories = []string{"bodywork", "dental", "deworming", "diet", "farrier", "fitting", "health", "lameness", "management", "other", "riding", "training", "vaccination", "vet"}

// Service contains business logic for care log management.
type Service struct {
	repo    *Repository
	subsSvc *subscriptions.Service
}

// NewService creates a care log service.
func NewService(repo *Repository, subsSvc *subscriptions.Service) *Service {
	return &Service{repo: repo, subsSvc: subsSvc}
}

// Input holds data for creating or updating a care log entry.
type Input struct {
	Date     string `json:"date"`
	Category string `json:"category"`
	Notes    string `json:"notes"`
}

// Validate checks the input for errors.
func (i *Input) Validate() validator.Errors {
	errs := validator.Errors{}
	validator.Required(errs, "date", i.Date)
	validator.OneOf(errs, "category", i.Category, validCategories)
	return errs
}

// GetByHorseID returns care logs for a horse, enforcing the feature gate.
func (s *Service) GetByHorseID(userID, horseID uuid.UUID) ([]models.CareLog, error) {
	if err := s.subsSvc.RequireFeature(userID, "care_logs"); err != nil {
		return nil, err
	}
	logs, err := s.repo.GetByHorseID(userID, horseID)
	if err != nil {
		return nil, err
	}
	if logs == nil {
		logs = []models.CareLog{}
	}
	return logs, nil
}

// Create adds a new care log entry.
func (s *Service) Create(userID, horseID uuid.UUID, input Input) (*models.CareLog, error) {
	if err := s.subsSvc.RequireFeature(userID, "care_logs"); err != nil {
		return nil, err
	}
	cl := &models.CareLog{
		UserID:   userID,
		HorseID:  horseID,
		Date:     input.Date,
		Category: input.Category,
		Notes:    input.Notes,
	}
	if err := s.repo.Create(cl); err != nil {
		return nil, err
	}
	return cl, nil
}

// Update modifies an existing care log entry.
func (s *Service) Update(userID, id uuid.UUID, input Input) (*models.CareLog, error) {
	if err := s.subsSvc.RequireFeature(userID, "care_logs"); err != nil {
		return nil, err
	}
	cl, err := s.repo.GetByID(userID, id)
	if err != nil {
		return nil, err
	}
	if cl == nil {
		return nil, nil
	}
	cl.Date = input.Date
	cl.Category = input.Category
	cl.Notes = input.Notes
	if err := s.repo.Update(cl); err != nil {
		return nil, err
	}
	return cl, nil
}

// Delete removes a care log entry.
func (s *Service) Delete(userID, id uuid.UUID) error {
	if err := s.subsSvc.RequireFeature(userID, "care_logs"); err != nil {
		return err
	}
	return s.repo.Delete(userID, id)
}
