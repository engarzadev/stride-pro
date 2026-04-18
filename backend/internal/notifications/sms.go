package notifications

import (
	"log"
)

// StubSMSSender logs SMS sends instead of actually sending them.
// Use this for development and testing.
type StubSMSSender struct{}

// NewStubSMSSender creates a stub SMS sender.
func NewStubSMSSender() *StubSMSSender {
	return &StubSMSSender{}
}

// Send logs the SMS details rather than sending a real message.
func (s *StubSMSSender) Send(recipient, subject, body string) error {
	log.Printf("[SMS STUB] To: %s | Body: %s", recipient, body)
	return nil
}

// TwilioSMSSender sends SMS messages via the Twilio API.
// TODO: Implement for v2 when Twilio integration is needed.
//
// type TwilioSMSSender struct {
// 	accountSID string
// 	authToken  string
// 	fromNumber string
// }
//
// func NewTwilioSMSSender(accountSID, authToken, fromNumber string) *TwilioSMSSender {
// 	return &TwilioSMSSender{
// 		accountSID: accountSID,
// 		authToken:  authToken,
// 		fromNumber: fromNumber,
// 	}
// }
//
// func (s *TwilioSMSSender) Send(recipient, subject, body string) error {
// 	// TODO: Use Twilio REST API to send SMS
// 	// 1. Build the request to https://api.twilio.com/2010-04-01/Accounts/{SID}/Messages.json
// 	// 2. Set Basic Auth with accountSID and authToken
// 	// 3. POST with To, From, Body parameters
// 	// 4. Handle response and errors
// 	return nil
// }
