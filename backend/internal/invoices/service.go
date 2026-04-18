package invoices

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	biz "github.com/stride-pro/backend/internal/business_settings"
	"github.com/stride-pro/backend/internal/models"
	"github.com/stride-pro/backend/internal/notifications"
	"github.com/stride-pro/backend/pkg/validator"
)

// Service contains business logic for invoice management.
type Service struct {
	repo        *Repository
	bizSettings *biz.Service
	emailSender notifications.HTMLSender
}

// NewService creates an invoice service.
func NewService(repo *Repository, bizSettings *biz.Service, emailSender notifications.HTMLSender) *Service {
	return &Service{repo: repo, bizSettings: bizSettings, emailSender: emailSender}
}

// ItemInput holds data for a single invoice line item.
type ItemInput struct {
	Description string  `json:"description"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	Notes       string  `json:"notes"`
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

// SendInvoice sets the invoice status to "sent" and emails it to the client.
func (s *Service) SendInvoice(userID, invoiceID uuid.UUID) error {
	inv, err := s.repo.GetByID(userID, invoiceID)
	if err != nil {
		return fmt.Errorf("fetching invoice: %w", err)
	}
	if inv == nil {
		return fmt.Errorf("invoice not found")
	}
	if inv.Client == nil || inv.Client.Email == "" {
		return fmt.Errorf("client has no email address")
	}

	bs, err := s.bizSettings.Get(userID)
	if err != nil {
		return fmt.Errorf("fetching business settings: %w", err)
	}

	html := buildInvoiceEmail(inv, bs)
	subject := fmt.Sprintf("Invoice from %s", bs.BusinessName)
	if bs.BusinessName == "" {
		subject = "Your Invoice"
	}

	if err := s.emailSender.SendHTML(inv.Client.Email, subject, html); err != nil {
		return fmt.Errorf("sending invoice email: %w", err)
	}

	return s.repo.UpdateStatus(userID, invoiceID, "sent")
}

// buildInvoiceEmail renders an HTML invoice email body.
func buildInvoiceEmail(inv *models.Invoice, bs *models.BusinessSettings) string {
	invoiceNum := inv.ID.String()[:8]

	itemRows := ""
	for _, item := range inv.Items {
		itemRows += fmt.Sprintf(`
			<tr>
				<td style="padding:8px;border-bottom:1px solid #eee;">%s</td>
				<td style="padding:8px;border-bottom:1px solid #eee;text-align:center;">%d</td>
				<td style="padding:8px;border-bottom:1px solid #eee;text-align:right;">$%.2f</td>
				<td style="padding:8px;border-bottom:1px solid #eee;text-align:right;">$%.2f</td>
			</tr>`, item.Description, item.Quantity, item.UnitPrice, item.Amount)
	}

	businessInfo := ""
	if bs.BusinessName != "" {
		businessInfo += fmt.Sprintf("<strong>%s</strong><br>", bs.BusinessName)
	}
	if bs.Phone != "" {
		businessInfo += fmt.Sprintf("%s<br>", bs.Phone)
	}
	if bs.Address != "" {
		businessInfo += fmt.Sprintf("%s<br>", bs.Address)
	}
	if bs.Email != "" {
		businessInfo += fmt.Sprintf("%s", bs.Email)
	}

	invoiceMessage := ""
	if bs.InvoiceMessage != "" {
		invoiceMessage = fmt.Sprintf(`<p style="margin-top:24px;color:#555;">%s</p>`, bs.InvoiceMessage)
	}

	clientName := ""
	if inv.Client != nil {
		clientName = inv.Client.FirstName + " " + inv.Client.LastName
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<body style="font-family:Arial,sans-serif;max-width:600px;margin:0 auto;padding:20px;color:#333;">
  <div style="border-bottom:2px solid #333;padding-bottom:16px;margin-bottom:24px;">
    <div style="font-size:12px;color:#777;">%s</div>
    <h2 style="margin:4px 0;">Invoice #%s</h2>
  </div>

  <p>Hi %s,</p>
  <p>Please find your invoice below.</p>

  <table style="width:100%%;border-collapse:collapse;margin:24px 0;">
    <thead>
      <tr style="background:#f5f5f5;">
        <th style="padding:8px;text-align:left;">Description</th>
        <th style="padding:8px;text-align:center;">Qty</th>
        <th style="padding:8px;text-align:right;">Unit Price</th>
        <th style="padding:8px;text-align:right;">Amount</th>
      </tr>
    </thead>
    <tbody>%s</tbody>
  </table>

  <div style="text-align:right;margin-top:8px;">
    <strong style="font-size:18px;">Total: $%.2f</strong>
  </div>

  <div style="margin-top:16px;font-size:13px;color:#777;">
    Due: %s
  </div>

  %s

  <div style="margin-top:40px;border-top:1px solid #eee;padding-top:16px;font-size:12px;color:#999;">
    %s
  </div>
</body>
</html>`,
		bs.BusinessName,
		invoiceNum,
		clientName,
		itemRows,
		inv.Total,
		inv.DueDate.Format("January 2, 2006"),
		invoiceMessage,
		businessInfo,
	)
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
			Notes:       inp.Notes,
		}
		total += amount
	}
	return items, total
}
