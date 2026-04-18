package barns

import (
	"github.com/google/uuid"

	"github.com/stride-pro/backend/internal/models"
	"github.com/stride-pro/backend/internal/subscriptions"
	"github.com/stride-pro/backend/pkg/validator"
)

// Service contains business logic for barn management.
type Service struct {
	repo    *Repository
	subsSvc *subscriptions.Service
}

// NewService creates a barn service.
func NewService(repo *Repository, subsSvc *subscriptions.Service) *Service {
	return &Service{repo: repo, subsSvc: subsSvc}
}

// CreateInput holds data for creating or updating a barn.
type CreateInput struct {
	Name        string `json:"name"`
	ContactName string `json:"contact_name"`
	Address     string `json:"address"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	Notes       string `json:"notes"`
}

// Validate checks the input for errors.
func (i *CreateInput) Validate() validator.Errors {
	errs := validator.Errors{}
	validator.Required(errs, "name", i.Name)
	validator.Email(errs, "email", i.Email)
	return errs
}

// Create validates and creates a new barn, requiring the barn_management feature.
func (s *Service) Create(userID uuid.UUID, input CreateInput) (*models.Barn, error) {
	if err := s.subsSvc.RequireFeature(userID, "barn_management"); err != nil {
		return nil, err
	}

	b := &models.Barn{
		UserID:      userID,
		Name:        input.Name,
		ContactName: input.ContactName,
		Address:     input.Address,
		Phone:       input.Phone,
		Email:       input.Email,
		Notes:       input.Notes,
	}
	if err := s.repo.Create(b); err != nil {
		return nil, err
	}
	return b, nil
}

// GetByID returns a barn by ID for the given user.
func (s *Service) GetByID(userID, barnID uuid.UUID) (*models.Barn, error) {
	return s.repo.GetByID(userID, barnID)
}

// GetAll returns all barns for the given user.
func (s *Service) GetAll(userID uuid.UUID) ([]models.Barn, error) {
	return s.repo.GetAllByUserID(userID)
}

// Update modifies an existing barn.
func (s *Service) Update(userID, barnID uuid.UUID, input CreateInput) (*models.Barn, error) {
	b := &models.Barn{
		ID:          barnID,
		UserID:      userID,
		Name:        input.Name,
		ContactName: input.ContactName,
		Address:     input.Address,
		Phone:       input.Phone,
		Email:       input.Email,
		Notes:       input.Notes,
	}
	if err := s.repo.Update(b); err != nil {
		return nil, err
	}
	return b, nil
}

// Delete removes a barn.
func (s *Service) Delete(userID, barnID uuid.UUID) error {
	return s.repo.Delete(userID, barnID)
}
