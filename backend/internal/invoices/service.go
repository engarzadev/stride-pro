package invoices

import (
	"time"

	"github.com/google/uuid"

	"github.com/stride-pro/backend/internal/models"
	"github.com/stride-pro/backend/pkg/validator"
)

// Service contains business logic for invoice management.
type Service struct {
	repo *Repository
}

// NewService creates an invoice service.
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// ItemInput holds data for a single invoice line item.
type ItemInput struct {
	Description string  `json:"description"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
}

// CreateInput holds data for creating or updating an invoice.
type CreateInput struct {
	ClientID uuid.UUID   `json:"client_id"`
	Status   string      `json:"status"`
	DueDate  time.Time   `json:"due_date"`
	Notes    string      `json:"notes"`
	Items    []ItemInput `json:"items"`
}

// Validate checks the input for errors.
func (i *CreateInput) Validate() validator.Errors {
	errs := validator.Errors{}
	validator.OneOf(errs, "status", i.Status, []string{"draft", "sent", "paid", "overdue", "cancelled"})
	if i.DueDate.IsZero() {
		errs["due_date"] = "due_date is required"
	}
	if len(i.Items) == 0 {
		errs["items"] = "at least one item is required"
	}
	for idx, item := range i.Items {
		if item.Description == "" {
			errs["items"] = "item description is required"
			_ = idx
			break
		}
		if item.Quantity < 1 {
			errs["items"] = "item quantity must be at least 1"
			break
		}
	}
	return errs
}

// Create validates and creates a new invoice, calculating totals from items.
func (s *Service) Create(userID uuid.UUID, input CreateInput) (*models.Invoice, error) {
	items, total := buildItems(input.Items)

	inv := &models.Invoice{
		UserID:   userID,
		ClientID: input.ClientID,
		Status:   input.Status,
		DueDate:  input.DueDate,
		Total:    total,
		Notes:    input.Notes,
		Items:    items,
	}
	if err := s.repo.Create(inv); err != nil {
		return nil, err
	}
	return inv, nil
}

// GetByID returns an invoice by ID for the given user (with items).
func (s *Service) GetByID(userID, invoiceID uuid.UUID) (*models.Invoice, error) {
	return s.repo.GetByID(userID, invoiceID)
}

// GetAll returns all invoices for the given user.
func (s *Service) GetAll(userID uuid.UUID) ([]models.Invoice, error) {
	return s.repo.GetAllByUserID(userID)
}

// GetByClientID returns all invoices for a specific client.
func (s *Service) GetByClientID(userID, clientID uuid.UUID) ([]models.Invoice, error) {
	return s.repo.GetByClientID(userID, clientID)
}

// Update modifies an existing invoice, recalculating totals.
func (s *Service) Update(userID, invoiceID uuid.UUID, input CreateInput) (*models.Invoice, error) {
	items, total := buildItems(input.Items)

	inv := &models.Invoice{
		ID:       invoiceID,
		UserID:   userID,
		ClientID: input.ClientID,
		Status:   input.Status,
		DueDate:  input.DueDate,
		Total:    total,
		Notes:    input.Notes,
		Items:    items,
	}
	if err := s.repo.Update(inv); err != nil {
		return nil, err
	}
	return inv, nil
}

// UpdateStatus changes the status of an invoice.
func (s *Service) UpdateStatus(userID, invoiceID uuid.UUID, status string) error {
	return s.repo.UpdateStatus(userID, invoiceID, status)
}

// Delete removes an invoice.
func (s *Service) Delete(userID, invoiceID uuid.UUID) error {
	return s.repo.Delete(userID, invoiceID)
}

// buildItems converts input items to model items and calculates the total.
func buildItems(inputs []ItemInput) ([]models.InvoiceItem, float64) {
	var total float64
	items := make([]models.InvoiceItem, len(inputs))
	for i, inp := range inputs {
		amount := float64(inp.Quantity) * inp.UnitPrice
		items[i] = models.InvoiceItem{
			Description: inp.Description,
			Quantity:    inp.Quantity,
			UnitPrice:   inp.UnitPrice,
			Amount:      amount,
		}
		total += amount
	}
	return items, total
}
