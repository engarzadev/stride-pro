package horses

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/stride-pro/backend/internal/models"
	"github.com/stride-pro/backend/internal/subscriptions"
	"github.com/stride-pro/backend/pkg/validator"
)

// Service contains business logic for horse management.
type Service struct {
	repo    *Repository
	subsSvc *subscriptions.Service
}

// NewService creates a horse service.
func NewService(repo *Repository, subsSvc *subscriptions.Service) *Service {
	return &Service{repo: repo, subsSvc: subsSvc}
}

// CreateInput holds data for creating or updating a horse.
type CreateInput struct {
	ClientID    uuid.UUID  `json:"client_id"`
	BarnID      *uuid.UUID `json:"barn_id"`
	Name        string     `json:"name"`
	Breed       string     `json:"breed"`
	Age         int        `json:"age"`
	Gender      string     `json:"gender"`
	Color       string     `json:"color"`
	Weight      float64    `json:"weight"`
	Notes       string     `json:"notes"`
	VetName     string     `json:"vet_name"`
	VetPhone    string     `json:"vet_phone"`
	FarrierName string     `json:"farrier_name"`
	FarrierPhone string    `json:"farrier_phone"`
}

// Validate checks the input for errors.
func (i *CreateInput) Validate() validator.Errors {
	errs := validator.Errors{}
	validator.Required(errs, "name", i.Name)
	validator.Required(errs, "breed", i.Breed)
	validator.MinValue(errs, "age", i.Age, 0)
	validator.OneOf(errs, "gender", i.Gender, []string{"stallion", "mare", "gelding"})
	return errs
}

// Create validates and creates a new horse, enforcing plan limits.
func (s *Service) Create(userID uuid.UUID, input CreateInput) (*models.Horse, error) {
	limit, err := s.subsSvc.GetHorseLimit(userID)
	if err != nil {
		return nil, fmt.Errorf("checking horse limit: %w", err)
	}
	if limit >= 0 {
		count, err := s.repo.CountByUserID(userID)
		if err != nil {
			return nil, fmt.Errorf("counting horses: %w", err)
		}
		if count >= limit {
			return nil, subscriptions.ErrLimitExceeded
		}
	}

	h := &models.Horse{
		UserID:       userID,
		ClientID:     input.ClientID,
		BarnID:       models.PtrToNullUUID(input.BarnID),
		Name:         input.Name,
		Breed:        input.Breed,
		Age:          input.Age,
		Gender:       input.Gender,
		Color:        input.Color,
		Weight:       input.Weight,
		Notes:        input.Notes,
		VetName:      input.VetName,
		VetPhone:     input.VetPhone,
		FarrierName:  input.FarrierName,
		FarrierPhone: input.FarrierPhone,
	}
	if err := s.repo.Create(h); err != nil {
		return nil, err
	}
	return h, nil
}

// GetByID returns a horse by ID for the given user.
func (s *Service) GetByID(userID, horseID uuid.UUID) (*models.Horse, error) {
	return s.repo.GetByID(userID, horseID)
}

// GetAll returns all horses for the given user.
func (s *Service) GetAll(userID uuid.UUID) ([]models.Horse, error) {
	return s.repo.GetAllByUserID(userID)
}

// GetByClientID returns all horses for a specific client.
func (s *Service) GetByClientID(userID, clientID uuid.UUID) ([]models.Horse, error) {
	return s.repo.GetByClientID(userID, clientID)
}

// GetByBarnID returns all horses at a specific barn.
func (s *Service) GetByBarnID(userID, barnID uuid.UUID) ([]models.Horse, error) {
	return s.repo.GetByBarnID(userID, barnID)
}

// Update modifies an existing horse.
func (s *Service) Update(userID, horseID uuid.UUID, input CreateInput) (*models.Horse, error) {
	h := &models.Horse{
		ID:           horseID,
		UserID:       userID,
		ClientID:     input.ClientID,
		BarnID:       models.PtrToNullUUID(input.BarnID),
		Name:         input.Name,
		Breed:        input.Breed,
		Age:          input.Age,
		Gender:       input.Gender,
		Color:        input.Color,
		Weight:       input.Weight,
		Notes:        input.Notes,
		VetName:      input.VetName,
		VetPhone:     input.VetPhone,
		FarrierName:  input.FarrierName,
		FarrierPhone: input.FarrierPhone,
	}
	if err := s.repo.Update(h); err != nil {
		return nil, err
	}
	return h, nil
}

// Delete removes a horse.
func (s *Service) Delete(userID, horseID uuid.UUID) error {
	return s.repo.Delete(userID, horseID)
}
