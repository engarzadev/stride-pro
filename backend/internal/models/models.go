// Package models defines the core domain types used throughout the application.
package models

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// User represents an authenticated user of the platform.
type User struct {
	ID               uuid.UUID `json:"id"`
	Email            string    `json:"email"`
	PasswordHash     string    `json:"-"`
	FirstName        string    `json:"first_name"`
	LastName         string    `json:"last_name"`
	Role             string    `json:"role"`
	SubscriptionTier string    `json:"subscription_tier"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// Client represents a horse owner or contact managed by a user.
type Client struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Address   string    `json:"address"`
	Notes     string    `json:"notes"`
	Horses    []*Horse  `json:"horses,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Horse represents a horse in the system, associated with a client and optionally a barn.
type Horse struct {
	ID        uuid.UUID      `json:"id"`
	UserID    uuid.UUID      `json:"user_id"`
	ClientID  uuid.NullUUID  `json:"client_id"`
	BarnID    uuid.NullUUID  `json:"barn_id"`
	Name      string         `json:"name"`
	Breed     string         `json:"breed"`
	Age       int            `json:"age"`
	Gender    string         `json:"gender"`
	Color     string         `json:"color"`
	Weight    float64        `json:"weight"`
	Notes       string         `json:"notes"`
	VetName     string         `json:"vet_name"`
	VetPhone    string         `json:"vet_phone"`
	FarrierName string         `json:"farrier_name"`
	FarrierPhone string        `json:"farrier_phone"`
	Client      *Client        `json:"client,omitempty"`
	Barn        *Barn          `json:"barn,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

// Barn represents a location where horses are kept.
type Barn struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	Name        string    `json:"name"`
	ContactName string    `json:"contact_name"`
	Address     string    `json:"address"`
	Phone       string    `json:"phone"`
	Email       string    `json:"email"`
	Notes       string    `json:"notes"`
	Horses      []*Horse  `json:"horses,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Appointment represents a scheduled visit or session.
type Appointment struct {
	ID          uuid.UUID     `json:"id"`
	UserID      uuid.UUID     `json:"user_id"`
	ClientID    uuid.UUID     `json:"client_id"`
	HorseID     uuid.UUID     `json:"horse_id"`
	BarnID      uuid.NullUUID `json:"barn_id"`
	ScheduledAt time.Time     `json:"-"`
	Date        string        `json:"date"`
	Time        string        `json:"time"`
	Duration    int           `json:"duration"`
	TravelTime  int           `json:"travel_time"`
	Status      string        `json:"status"`
	Type        string        `json:"type"`
	Notes       string        `json:"notes"`
	Client      *Client       `json:"client,omitempty"`
	Horse       *Horse        `json:"horse,omitempty"`
	Barn        *Barn         `json:"barn,omitempty"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

// Session records the details of a care session performed during an appointment.
type Session struct {
	ID              uuid.UUID       `json:"id"`
	UserID          uuid.UUID       `json:"user_id"`
	AppointmentID   uuid.UUID       `json:"appointment_id"`
	Type            string          `json:"type"`
	BodyZones       json.RawMessage `json:"body_zones"`
	Notes           string          `json:"notes"`
	Findings        string          `json:"findings"`
	Recommendations string          `json:"recommendations"`
	Appointment     *Appointment    `json:"appointment,omitempty"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

// Invoice represents a billing document sent to a client.
type Invoice struct {
	ID        uuid.UUID     `json:"id"`
	UserID    uuid.UUID     `json:"user_id"`
	ClientID  uuid.UUID     `json:"client_id"`
	Status    string        `json:"status"`
	DueDate   time.Time     `json:"due_date"`
	Total     float64       `json:"total"`
	Notes     string        `json:"notes"`
	Client    *Client       `json:"client,omitempty"`
	Items     []InvoiceItem `json:"items,omitempty"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

// InvoiceItem is a single line item on an invoice.
type InvoiceItem struct {
	ID          uuid.UUID `json:"id"`
	InvoiceID   uuid.UUID `json:"invoice_id"`
	Description string    `json:"description"`
	Quantity    int       `json:"quantity"`
	UnitPrice   float64   `json:"unit_price"`
	Amount      float64   `json:"amount"`
	Notes       string    `json:"notes,omitempty"`
}

// BusinessSettings holds a user's business profile used for invoicing.
type BusinessSettings struct {
	ID             uuid.UUID `json:"id"`
	UserID         uuid.UUID `json:"user_id"`
	BusinessName   string    `json:"business_name"`
	Email          string    `json:"email"`
	Phone          string    `json:"phone"`
	Address        string    `json:"address"`
	InvoiceMessage string    `json:"invoice_message"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// ServiceItem is a reusable service in the user's price catalog.
type ServiceItem struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID `json:"user_id"`
	Name         string    `json:"name"`
	DefaultPrice float64   `json:"default_price"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Notification represents a message sent to a user or client.
type Notification struct {
	ID        uuid.UUID    `json:"id"`
	UserID    uuid.UUID    `json:"user_id"`
	Type      string       `json:"type"`
	Channel   string       `json:"channel"`
	Recipient string       `json:"recipient"`
	Subject   string       `json:"subject"`
	Body      string       `json:"body"`
	Status    string       `json:"status"`
	SentAt    sql.NullTime `json:"sent_at"`
	CreatedAt time.Time    `json:"created_at"`
}

// Reminder represents an upcoming health or maintenance task for a horse.
type Reminder struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	HorseID    uuid.UUID `json:"horse_id"`
	Title      string    `json:"title"`
	DueDate    string    `json:"due_date"`
	Category   string    `json:"category"`
	Source     string    `json:"source"`
	IsComplete bool      `json:"is_complete"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// CareLog records a care event for a horse (e.g. farrier, vet, diet change).
type CareLog struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	HorseID   uuid.UUID `json:"horse_id"`
	Date      string    `json:"date"`
	Category  string    `json:"category"`
	Notes     string    `json:"notes"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Subscription represents a user's billing plan.
type Subscription struct {
	ID        uuid.UUID       `json:"id"`
	UserID    uuid.UUID       `json:"user_id"`
	Plan      string          `json:"plan"`
	Status    string          `json:"status"`
	Features  json.RawMessage `json:"features"`
	StartsAt  time.Time       `json:"starts_at"`
	EndsAt    time.Time       `json:"ends_at"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

// NullUUIDToPtr converts a sql NullUUID to a *uuid.UUID pointer.
func NullUUIDToPtr(n uuid.NullUUID) *uuid.UUID {
	if !n.Valid {
		return nil
	}
	return &n.UUID
}

// PtrToNullUUID converts a *uuid.UUID pointer to a sql NullUUID.
func PtrToNullUUID(u *uuid.UUID) uuid.NullUUID {
	if u == nil {
		return uuid.NullUUID{}
	}
	return uuid.NullUUID{UUID: *u, Valid: true}
}
