// Package config handles application configuration loaded from environment variables.
package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Config holds all application configuration values.
type Config struct {
	DatabaseURL string
	JWTSecret   string
	Port        string
	Environment string

	// External service keys (for future use)
	SendGridAPIKey    string
	TwilioAccountSID  string
	TwilioAuthToken   string
	TwilioPhoneNumber string
}

// Load reads configuration from environment variables.
// It attempts to load a .env file first but does not fail if one is not found.
func Load() (*Config, error) {
	// Best-effort load of .env file
	_ = godotenv.Load()

	cfg := &Config{
		DatabaseURL:       getEnv("DATABASE_URL", ""),
		JWTSecret:         getEnv("JWT_SECRET", ""),
		Port:              getEnv("PORT", "8080"),
		Environment:       getEnv("ENVIRONMENT", "dev"),
		SendGridAPIKey:    getEnv("SENDGRID_API_KEY", ""),
		TwilioAccountSID:  getEnv("TWILIO_ACCOUNT_SID", ""),
		TwilioAuthToken:   getEnv("TWILIO_AUTH_TOKEN", ""),
		TwilioPhoneNumber: getEnv("TWILIO_PHONE_NUMBER", ""),
	}

	return cfg, nil
}

// IsProd returns true if the environment is production.
func (c *Config) IsProd() bool {
	return c.Environment == "prod"
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
