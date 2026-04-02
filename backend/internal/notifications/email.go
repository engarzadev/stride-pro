package notifications

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// HTMLSender is the interface for sending HTML-formatted email.
type HTMLSender interface {
	SendHTML(recipient, subject, htmlBody string) error
}

// EmailSender combines plain-text and HTML email sending.
type EmailSender interface {
	Sender
	HTMLSender
}

// StubEmailSender logs email sends instead of actually sending them.
// Use this for development and testing.
type StubEmailSender struct{}

// NewStubEmailSender creates a stub email sender.
func NewStubEmailSender() *StubEmailSender {
	return &StubEmailSender{}
}

// Send logs the email details rather than sending a real email.
func (s *StubEmailSender) Send(recipient, subject, body string) error {
	log.Printf("[EMAIL STUB] To: %s | Subject: %s | Body: %s", recipient, subject, body)
	return nil
}

// SendHTML logs the HTML email details rather than sending a real email.
func (s *StubEmailSender) SendHTML(recipient, subject, htmlBody string) error {
	log.Printf("[EMAIL STUB HTML] To: %s | Subject: %s | Body: %s", recipient, subject, htmlBody)
	return nil
}

// SendGridEmailSender sends emails via the SendGrid v3 API.
type SendGridEmailSender struct {
	apiKey    string
	fromEmail string
	fromName  string
	client    *http.Client
}

// NewSendGridEmailSender creates a SendGrid email sender.
func NewSendGridEmailSender(apiKey, fromEmail, fromName string) *SendGridEmailSender {
	return &SendGridEmailSender{
		apiKey:    apiKey,
		fromEmail: fromEmail,
		fromName:  fromName,
		client:    &http.Client{},
	}
}

// Send delivers a plain-text email via the SendGrid v3 mail/send endpoint.
func (s *SendGridEmailSender) Send(recipient, subject, body string) error {
	return s.send(recipient, subject, "text/plain", body)
}

// SendHTML delivers an HTML email via the SendGrid v3 mail/send endpoint.
func (s *SendGridEmailSender) SendHTML(recipient, subject, htmlBody string) error {
	return s.send(recipient, subject, "text/html", htmlBody)
}

func (s *SendGridEmailSender) send(recipient, subject, contentType, body string) error {
	payload := map[string]any{
		"personalizations": []map[string]any{
			{"to": []map[string]string{{"email": recipient}}},
		},
		"from": map[string]string{
			"email": s.fromEmail,
			"name":  s.fromName,
		},
		"subject": subject,
		"content": []map[string]string{
			{"type": contentType, "value": body},
		},
	}

	b, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshalling sendgrid payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, "https://api.sendgrid.com/v3/mail/send", bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("building sendgrid request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("sending sendgrid request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("sendgrid returned %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}
