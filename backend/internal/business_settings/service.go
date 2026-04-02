package business_settings

import (
	"github.com/google/uuid"

	"github.com/stride-pro/backend/internal/models"
)

// Service contains business logic for business settings management.
type Service struct {
	repo *Repository
}

// NewService creates a business settings service.
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// Get returns the business settings for a user.
func (s *Service) Get(userID uuid.UUID) (*models.BusinessSettings, error) {
	bs, err := s.repo.Get(userID)
	if err != nil {
		return nil, err
	}
	if bs == nil {
		// Return empty defaults so the frontend always gets a valid object
		return &models.BusinessSettings{UserID: userID}, nil
	}
	return bs, nil
}

// Upsert saves business settings for a user.
func (s *Service) Upsert(userID uuid.UUID, input UpsertInput) (*models.BusinessSettings, error) {
	bs := &models.BusinessSettings{
		UserID:         userID,
		BusinessName:   input.BusinessName,
		Email:          input.Email,
		Phone:          input.Phone,
		Address:        input.Address,
		InvoiceMessage: input.InvoiceMessage,
	}
	if err := s.repo.Upsert(bs); err != nil {
		return nil, err
	}
	return bs, nil
}

// UpsertInput holds data for creating or updating business settings.
type UpsertInput struct {
	BusinessName   string `json:"business_name"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	Address        string `json:"address"`
	InvoiceMessage string `json:"invoice_message"`
}
