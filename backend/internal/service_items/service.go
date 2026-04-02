package service_items

import (
	"github.com/google/uuid"

	"github.com/stride-pro/backend/internal/models"
)

// Service contains business logic for service item management.
type Service struct {
	repo *Repository
}

// NewService creates a service item service.
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// ItemInput holds data for creating or updating a service item.
type ItemInput struct {
	Name         string  `json:"name"`
	DefaultPrice float64 `json:"default_price"`
}

// GetAll returns all service items for a user.
func (s *Service) GetAll(userID uuid.UUID) ([]models.ServiceItem, error) {
	return s.repo.GetAll(userID)
}

// Create adds a new service item.
func (s *Service) Create(userID uuid.UUID, input ItemInput) (*models.ServiceItem, error) {
	item := &models.ServiceItem{
		UserID:       userID,
		Name:         input.Name,
		DefaultPrice: input.DefaultPrice,
	}
	if err := s.repo.Create(item); err != nil {
		return nil, err
	}
	return item, nil
}

// Update modifies an existing service item.
func (s *Service) Update(userID, itemID uuid.UUID, input ItemInput) (*models.ServiceItem, error) {
	item := &models.ServiceItem{
		ID:           itemID,
		UserID:       userID,
		Name:         input.Name,
		DefaultPrice: input.DefaultPrice,
	}
	if err := s.repo.Update(item); err != nil {
		return nil, err
	}
	return item, nil
}

// Delete removes a service item.
func (s *Service) Delete(userID, itemID uuid.UUID) error {
	return s.repo.Delete(userID, itemID)
}
