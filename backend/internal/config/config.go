// Package config handles application configuration loaded from environment variables.
package config

import (
	"errors"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds all application configuration values.
type Config struct {
	DatabaseURL string
	JWTSecret   string
	Port        string
	Environment string

	// CORS — comma-separated list of allowed origins, e.g. "https://app.example.com"
	AllowedOrigins []string

	// TLS (direct) — if both are set the server terminates TLS itself.
	// An HTTP→HTTPS redirect server will also start on HTTPPort.
	TLSCertFile string
	TLSKeyFile  string

	// HTTPPort is the port the HTTP redirect server listens on when direct TLS
	// is enabled. Defaults to 80. Not used in proxy mode.
	HTTPPort string

	// TLSProxyMode — set to true when TLS is terminated by an upstream proxy
	// (e.g. AWS ALB, Cloudflare, nginx). The app trusts X-Forwarded-Proto and
	// redirects any request where it is not "https".
	TLSProxyMode bool

	// External service keys
	SendGridAPIKey    string
	SendGridFromEmail string
	SendGridFromName  string
	TwilioAccountSID  string
	TwilioAuthToken   string
	TwilioPhoneNumber string

	// AppBaseURL is the public-facing frontend URL used to build email links
	// (e.g. password reset). Example: "https://app.stridepro.com"
	AppBaseURL string
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
		TLSCertFile:       getEnv("TLS_CERT_FILE", ""),
		TLSKeyFile:        getEnv("TLS_KEY_FILE", ""),
		HTTPPort:          getEnv("HTTP_PORT", "80"),
		TLSProxyMode:      getEnv("TLS_PROXY_MODE", "") == "true",
		SendGridAPIKey:    getEnv("SENDGRID_API_KEY", ""),
		SendGridFromEmail: getEnv("SENDGRID_FROM_EMAIL", ""),
		SendGridFromName:  getEnv("SENDGRID_FROM_NAME", "Stride Pro"),
		TwilioAccountSID:  getEnv("TWILIO_ACCOUNT_SID", ""),
		TwilioAuthToken:   getEnv("TWILIO_AUTH_TOKEN", ""),
		TwilioPhoneNumber: getEnv("TWILIO_PHONE_NUMBER", ""),
		AppBaseURL:        getEnv("APP_BASE_URL", "http://localhost:4200"),
	}

	// Parse allowed origins from comma-separated env var
	rawOrigins := getEnv("ALLOWED_ORIGINS", "http://localhost:4200")
	for _, o := range strings.Split(rawOrigins, ",") {
		if trimmed := strings.TrimSpace(o); trimmed != "" {
			cfg.AllowedOrigins = append(cfg.AllowedOrigins, trimmed)
		}
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// validate checks that required fields are present and meet minimum security requirements.
func (c *Config) validate() error {
	if c.DatabaseURL == "" {
		return errors.New("DATABASE_URL is required")
	}
	if c.JWTSecret == "" {
		return errors.New("JWT_SECRET is required")
	}
	// Enforce a minimum secret length to prevent weak secrets in any environment.
	// 32 characters provides at least 256 bits of entropy for an ASCII secret.
	if len(c.JWTSecret) < 32 {
		return errors.New("JWT_SECRET must be at least 32 characters long")
	}
	// In production, HTTPS must be explicitly configured — either directly via
	// TLS cert/key files, or via a trusted upstream proxy. Starting without any
	// HTTPS coverage is not allowed.
	if c.IsProd() && !c.TLSEnabled() && !c.TLSProxyMode {
		return errors.New(
			"HTTPS is required in production: " +
				"set TLS_CERT_FILE + TLS_KEY_FILE for direct TLS, " +
				"or set TLS_PROXY_MODE=true if TLS is terminated by an upstream proxy",
		)
	}
	return nil
}

// IsProd returns true if the environment is production.
func (c *Config) IsProd() bool {
	return c.Environment == "prod" || c.Environment == "production"
}

// TLSEnabled returns true when both TLS certificate and key paths are configured.
func (c *Config) TLSEnabled() bool {
	return c.TLSCertFile != "" && c.TLSKeyFile != ""
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
