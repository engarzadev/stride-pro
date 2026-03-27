package notifications

import (
	"log"
)

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

// SendGridEmailSender sends emails via the SendGrid API.
// TODO: Implement for v2 when SendGrid integration is needed.
//
// type SendGridEmailSender struct {
// 	apiKey    string
// 	fromEmail string
// 	fromName  string
// }
//
// func NewSendGridEmailSender(apiKey, fromEmail, fromName string) *SendGridEmailSender {
// 	return &SendGridEmailSender{
// 		apiKey:    apiKey,
// 		fromEmail: fromEmail,
// 		fromName:  fromName,
// 	}
// }
//
// func (s *SendGridEmailSender) Send(recipient, subject, body string) error {
// 	// TODO: Use SendGrid v3 API to send email
// 	// 1. Build the mail message
// 	// 2. POST to https://api.sendgrid.com/v3/mail/send
// 	// 3. Handle response and errors
// 	return nil
// }
