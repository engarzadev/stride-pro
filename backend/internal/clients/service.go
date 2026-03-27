package clients

import (
	"github.com/google/uuid"

	"github.com/stride-pro/backend/internal/models"
	"github.com/stride-pro/backend/pkg/validator"
)

// Service contains business logic for client management.
type Service struct {
	repo *Repository
}

// NewService creates a client service.
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// CreateInput holds data for creating a client.
type CreateInput struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
	Notes     string `json:"notes"`
}

// Validate checks the create input for errors.
func (i *CreateInput) Validate() validator.Errors {
	errs := validator.Errors{}
	validator.Required(errs, "first_name", i.FirstName)
	validator.Required(errs, "last_name", i.LastName)
	validator.Email(errs, "email", i.Email)
	return errs
}

// Create validates and creates a new client.
func (s *Service) Create(userID uuid.UUID, input CreateInput) (*models.Client, error) {
	c := &models.Client{
		UserID:    userID,
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		Phone:     input.Phone,
		Address:   input.Address,
		Notes:     input.Notes,
	}
	if err := s.repo.Create(c); err != nil {
		return nil, err
	}
	return c, nil
}

// GetByID returns a client by ID for the given user.
func (s *Service) GetByID(userID, clientID uuid.UUID) (*models.Client, error) {
	return s.repo.GetByID(userID, clientID)
}

// GetAll returns all clients for the given user.
func (s *Service) GetAll(userID uuid.UUID) ([]models.Client, error) {
	return s.repo.GetAllByUserID(userID)
}

// Update modifies an existing client.
func (s *Service) Update(userID, clientID uuid.UUID, input CreateInput) (*models.Client, error) {
	c := &models.Client{
		ID:        clientID,
		UserID:    userID,
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		Phone:     input.Phone,
		Address:   input.Address,
		Notes:     input.Notes,
	}
	if err := s.repo.Update(c); err != nil {
		return nil, err
	}
	return c, nil
}

// Delete removes a client.
func (s *Service) Delete(userID, clientID uuid.UUID) error {
	return s.repo.Delete(userID, clientID)
}
