// Package invoices manages invoice data access and business logic.
package invoices

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/stride-pro/backend/internal/database"
	"github.com/stride-pro/backend/internal/models"
)

// Repository handles invoice persistence.
type Repository struct {
	db *database.DB
}

// NewRepository creates an invoice repository.
func NewRepository(db *database.DB) *Repository {
	return &Repository{db: db}
}

const invoiceSelectCols = `i.id, i.user_id, i.client_id, i.status, i.due_date, i.total, i.notes, i.created_at, i.updated_at, c.id, c.first_name, c.last_name`
const invoiceJoins = `FROM invoices i LEFT JOIN clients c ON i.client_id = c.id`

func scanInvoice(scanner interface{ Scan(...interface{}) error }) (*models.Invoice, error) {
	inv := &models.Invoice{}
	var clientID uuid.NullUUID
	var clientFirstName, clientLastName sql.NullString
	err := scanner.Scan(
		&inv.ID, &inv.UserID, &inv.ClientID, &inv.Status,
		&inv.DueDate, &inv.Total, &inv.Notes, &inv.CreatedAt, &inv.UpdatedAt,
		&clientID, &clientFirstName, &clientLastName,
	)
	if err != nil {
		return nil, err
	}
	if clientID.Valid {
		inv.Client = &models.Client{
			ID:        clientID.UUID,
			FirstName: clientFirstName.String,
			LastName:  clientLastName.String,
		}
	}
	return inv, nil
}

// Create inserts a new invoice and its items in a transaction.
func (r *Repository) Create(inv *models.Invoice) error {
	inv.ID = uuid.New()
	inv.CreatedAt = time.Now()
	inv.UpdatedAt = time.Now()

	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("starting transaction: %w", err)
	}

	_, err = tx.Exec(`
		INSERT INTO invoices (id, user_id, client_id, status, due_date, total, notes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		inv.ID, inv.UserID, inv.ClientID, inv.Status,
		inv.DueDate, inv.Total, inv.Notes, inv.CreatedAt, inv.UpdatedAt,
	)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("inserting invoice: %w", err)
	}

	for i := range inv.Items {
		inv.Items[i].ID = uuid.New()
		inv.Items[i].InvoiceID = inv.ID
		_, err = tx.Exec(`
			INSERT INTO invoice_items (id, invoice_id, description, quantity, unit_price, amount)
			VALUES ($1, $2, $3, $4, $5, $6)`,
			inv.Items[i].ID, inv.Items[i].InvoiceID, inv.Items[i].Description,
			inv.Items[i].Quantity, inv.Items[i].UnitPrice, inv.Items[i].Amount,
		)
		if err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("inserting invoice item: %w", err)
		}
	}

	return tx.Commit()
}

// GetByID returns a single invoice with its items, scoped to the user.
func (r *Repository) GetByID(userID, invoiceID uuid.UUID) (*models.Invoice, error) {
	inv, err := scanInvoice(r.db.QueryRow(
		`SELECT `+invoiceSelectCols+` `+invoiceJoins+` WHERE i.id = $1 AND i.user_id = $2`,
		invoiceID, userID,
	))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("querying invoice: %w", err)
	}

	items, err := r.getItems(invoiceID)
	if err != nil {
		return nil, err
	}
	inv.Items = items

	return inv, nil
}

// GetAllByUserID returns all invoices belonging to a user (without items for performance).
func (r *Repository) GetAllByUserID(userID uuid.UUID) ([]models.Invoice, error) {
	rows, err := r.db.Query(
		`SELECT `+invoiceSelectCols+` `+invoiceJoins+` WHERE i.user_id = $1 ORDER BY i.created_at DESC`, userID,
	)
	if err != nil {
		return nil, fmt.Errorf("querying invoices: %w", err)
	}
	defer rows.Close()

	var invoices []models.Invoice
	for rows.Next() {
		inv, err := scanInvoice(rows)
		if err != nil {
			return nil, fmt.Errorf("scanning invoice: %w", err)
		}
		invoices = append(invoices, *inv)
	}
	return invoices, rows.Err()
}

// GetByClientID returns all invoices for a specific client.
func (r *Repository) GetByClientID(userID, clientID uuid.UUID) ([]models.Invoice, error) {
	rows, err := r.db.Query(
		`SELECT `+invoiceSelectCols+` `+invoiceJoins+` WHERE i.user_id = $1 AND i.client_id = $2 ORDER BY i.created_at DESC`,
		userID, clientID,
	)
	if err != nil {
		return nil, fmt.Errorf("querying invoices by client: %w", err)
	}
	defer rows.Close()

	var invoices []models.Invoice
	for rows.Next() {
		inv, err := scanInvoice(rows)
		if err != nil {
			return nil, fmt.Errorf("scanning invoice: %w", err)
		}
		invoices = append(invoices, *inv)
	}
	return invoices, rows.Err()
}

// Update modifies an existing invoice and replaces its items.
func (r *Repository) Update(inv *models.Invoice) error {
	inv.UpdatedAt = time.Now()

	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("starting transaction: %w", err)
	}

	result, err := tx.Exec(`
		UPDATE invoices SET client_id=$1, status=$2, due_date=$3, total=$4, notes=$5, updated_at=$6
		WHERE id=$7 AND user_id=$8`,
		inv.ClientID, inv.Status, inv.DueDate, inv.Total, inv.Notes, inv.UpdatedAt, inv.ID, inv.UserID,
	)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("updating invoice: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		_ = tx.Rollback()
		return fmt.Errorf("invoice not found")
	}

	// Replace all items
	if _, err := tx.Exec("DELETE FROM invoice_items WHERE invoice_id = $1", inv.ID); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("deleting old invoice items: %w", err)
	}

	for i := range inv.Items {
		inv.Items[i].ID = uuid.New()
		inv.Items[i].InvoiceID = inv.ID
		_, err = tx.Exec(`
			INSERT INTO invoice_items (id, invoice_id, description, quantity, unit_price, amount)
			VALUES ($1, $2, $3, $4, $5, $6)`,
			inv.Items[i].ID, inv.Items[i].InvoiceID, inv.Items[i].Description,
			inv.Items[i].Quantity, inv.Items[i].UnitPrice, inv.Items[i].Amount,
		)
		if err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("inserting invoice item: %w", err)
		}
	}

	return tx.Commit()
}

// UpdateStatus changes the status of an invoice.
func (r *Repository) UpdateStatus(userID, invoiceID uuid.UUID, status string) error {
	result, err := r.db.Exec(
		`UPDATE invoices SET status = $1, updated_at = $2 WHERE id = $3 AND user_id = $4`,
		status, time.Now(), invoiceID, userID,
	)
	if err != nil {
		return fmt.Errorf("updating invoice status: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("invoice not found")
	}
	return nil
}

// Delete removes an invoice and its items.
func (r *Repository) Delete(userID, invoiceID uuid.UUID) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("starting transaction: %w", err)
	}

	if _, err := tx.Exec("DELETE FROM invoice_items WHERE invoice_id = $1", invoiceID); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("deleting invoice items: %w", err)
	}

	result, err := tx.Exec("DELETE FROM invoices WHERE id = $1 AND user_id = $2", invoiceID, userID)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("deleting invoice: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		_ = tx.Rollback()
		return fmt.Errorf("invoice not found")
	}

	return tx.Commit()
}

func (r *Repository) getItems(invoiceID uuid.UUID) ([]models.InvoiceItem, error) {
	rows, err := r.db.Query(`
		SELECT id, invoice_id, description, quantity, unit_price, amount
		FROM invoice_items WHERE invoice_id = $1 ORDER BY id`,
		invoiceID,
	)
	if err != nil {
		return nil, fmt.Errorf("querying invoice items: %w", err)
	}
	defer rows.Close()

	var items []models.InvoiceItem
	for rows.Next() {
		var item models.InvoiceItem
		if err := rows.Scan(&item.ID, &item.InvoiceID, &item.Description, &item.Quantity, &item.UnitPrice, &item.Amount); err != nil {
			return nil, fmt.Errorf("scanning invoice item: %w", err)
		}
		items = append(items, item)
	}
	return items, rows.Err()
}
