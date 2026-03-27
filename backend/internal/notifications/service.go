// Package notifications handles sending notifications via email and SMS.
package notifications

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"

	"github.com/stride-pro/backend/internal/database"
	"github.com/stride-pro/backend/internal/models"
)

// Sender is the interface for sending a notification through a specific channel.
type Sender interface {
	Send(recipient, subject, body string) error
}

// Service manages notification dispatch and persistence.
type Service struct {
	db    *database.DB
	email Sender
	sms   Sender
}

// NewService creates a notification service with the given channel senders.
func NewService(db *database.DB, email, sms Sender) *Service {
	return &Service{
		db:    db,
		email: email,
		sms:   sms,
	}
}

// Send dispatches a notification through the appropriate channel and records the result.
func (s *Service) Send(n *models.Notification) error {
	n.ID = uuid.New()
	n.Status = "pending"
	n.CreatedAt = time.Now()

	var sender Sender
	switch n.Channel {
	case "email":
		sender = s.email
	case "sms":
		sender = s.sms
	default:
		return fmt.Errorf("unsupported notification channel: %s", n.Channel)
	}

	if err := sender.Send(n.Recipient, n.Subject, n.Body); err != nil {
		n.Status = "failed"
		s.persist(n)
		return fmt.Errorf("sending notification: %w", err)
	}

	n.Status = "sent"
	now := time.Now()
	n.SentAt.Time = now
	n.SentAt.Valid = true
	s.persist(n)

	return nil
}

func (s *Service) persist(n *models.Notification) {
	_, err := s.db.Exec(`
		INSERT INTO notifications (id, user_id, type, channel, recipient, subject, body, status, sent_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		n.ID, n.UserID, n.Type, n.Channel, n.Recipient, n.Subject, n.Body, n.Status, n.SentAt, n.CreatedAt,
	)
	if err != nil {
		log.Printf("failed to persist notification %s: %v", n.ID, err)
	}
}
