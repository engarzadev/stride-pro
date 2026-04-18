// Package subscriptions defines subscription plans and feature flags.
package subscriptions

// Plan defines a subscription plan with its available features.
type Plan struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       float64  `json:"price"`
	Features    []string `json:"features"`
}

// Plans holds all available subscription plans.
var Plans = map[string]Plan{
	"free": {
		ID:          "free",
		Name:        "Free",
		Description: "Basic access with limited features",
		Price:       0,
		Features: []string{
			"clients_max_10",
			"horses_max_20",
			"appointments_basic",
			"invoices_basic",
			"care_log_reminders",
		},
	},
	"base": {
		ID:          "base",
		Name:        "Base",
		Description: "Full access for independent practitioners",
		Price:       29.99,
		Features: []string{
			"clients_unlimited",
			"horses_unlimited",
			"appointments_full",
			"invoices_full",
			"session_notes",
			"barn_management",
			"email_notifications",
			"care_logs",
		},
	},
	"trainer_addon": {
		ID:          "trainer_addon",
		Name:        "Trainer Add-on",
		Description: "Extended features for trainers and larger practices",
		Price:       49.99,
		Features: []string{
			"clients_unlimited",
			"horses_unlimited",
			"appointments_full",
			"invoices_full",
			"session_notes",
			"barn_management",
			"email_notifications",
			"sms_notifications",
			"multi_horse_sessions",
			"advanced_reporting",
			"client_portal",
			"care_logs",
			"care_log_reminders",
		},
	},
	"enterprise": {
		ID:          "enterprise",
		Name:        "Enterprise",
		Description: "Full platform access for large organizations",
		Price:       99.99,
		Features: []string{
			"clients_unlimited",
			"horses_unlimited",
			"appointments_full",
			"invoices_full",
			"session_notes",
			"barn_management",
			"email_notifications",
			"sms_notifications",
			"multi_horse_sessions",
			"advanced_reporting",
			"client_portal",
			"api_access",
			"custom_branding",
			"priority_support",
			"care_logs",
			"care_log_reminders",
		},
	},
}

// FeatureDescriptions maps feature flags to human-readable descriptions.
var FeatureDescriptions = map[string]string{
	"clients_max_10":       "Up to 10 clients",
	"clients_unlimited":    "Unlimited clients",
	"horses_max_20":        "Up to 20 horses",
	"horses_unlimited":     "Unlimited horses",
	"appointments_basic":   "Basic scheduling",
	"appointments_full":    "Full scheduling with reminders",
	"invoices_basic":       "Basic invoicing",
	"invoices_full":        "Full invoicing with templates",
	"session_notes":        "Detailed session notes and findings",
	"barn_management":      "Barn/location management",
	"email_notifications":  "Email notifications",
	"sms_notifications":    "SMS notifications",
	"multi_horse_sessions": "Multi-horse session tracking",
	"advanced_reporting":   "Advanced reporting and analytics",
	"client_portal":        "Client self-service portal",
	"api_access":           "REST API access",
	"custom_branding":      "Custom branding",
	"priority_support":     "Priority support",
	"care_logs":            "Horse care log (farrier, vet, diet, and more)",
	"care_log_reminders":   "Automated reminders for recurring care",
}
