package notifications

import (
	"fmt"
	"time"
)

// Template defines a notification template with subject and body generators.
type Template struct {
	Subject func(data map[string]string) string
	Body    func(data map[string]string) string
}

// Templates contains pre-defined notification templates keyed by name.
var Templates = map[string]Template{
	"appointment_reminder": {
		Subject: func(data map[string]string) string {
			return fmt.Sprintf("Appointment Reminder - %s", data["date"])
		},
		Body: func(data map[string]string) string {
			return fmt.Sprintf(
				"Hi %s,\n\nThis is a reminder that you have an appointment scheduled for %s at %s for %s.\n\nPlease let us know if you need to reschedule.\n\nBest regards,\n%s",
				data["client_name"],
				data["date"],
				data["time"],
				data["horse_name"],
				data["provider_name"],
			)
		},
	},
	"invoice_reminder": {
		Subject: func(data map[string]string) string {
			return fmt.Sprintf("Invoice #%s - Payment Reminder", data["invoice_number"])
		},
		Body: func(data map[string]string) string {
			return fmt.Sprintf(
				"Hi %s,\n\nThis is a friendly reminder that Invoice #%s for $%s is due on %s.\n\nPlease arrange payment at your earliest convenience.\n\nThank you,\n%s",
				data["client_name"],
				data["invoice_number"],
				data["amount"],
				data["due_date"],
				data["provider_name"],
			)
		},
	},
	"payment_confirmation": {
		Subject: func(_ map[string]string) string {
			return "Payment Received - Thank You!"
		},
		Body: func(data map[string]string) string {
			return fmt.Sprintf(
				"Hi %s,\n\nWe have received your payment of $%s for Invoice #%s.\n\nThank you for your prompt payment!\n\nBest regards,\n%s",
				data["client_name"],
				data["amount"],
				data["invoice_number"],
				data["provider_name"],
			)
		},
	},
	"booking_confirmation": {
		Subject: func(data map[string]string) string {
			return fmt.Sprintf("Booking Confirmed - %s", data["date"])
		},
		Body: func(data map[string]string) string {
			return fmt.Sprintf(
				"Hi %s,\n\nYour appointment has been confirmed for %s at %s.\n\nHorse: %s\nType: %s\nDuration: %s minutes\n\nWe look forward to seeing you!\n\nBest regards,\n%s",
				data["client_name"],
				data["date"],
				data["time"],
				data["horse_name"],
				data["appointment_type"],
				data["duration"],
				data["provider_name"],
			)
		},
	},
}

// RenderTemplate renders a notification template by name with the given data.
func RenderTemplate(name string, data map[string]string) (subject, body string, err error) {
	tmpl, ok := Templates[name]
	if !ok {
		return "", "", fmt.Errorf("template not found: %s", name)
	}

	// Ensure date defaults
	if _, ok := data["date"]; !ok {
		data["date"] = time.Now().Format("January 2, 2006")
	}

	return tmpl.Subject(data), tmpl.Body(data), nil
}
