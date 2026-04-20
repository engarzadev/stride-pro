// Package validator provides simple field validation helpers.
package validator

import (
	"fmt"
	"net/mail"
	"strings"
	"unicode"
)

// Errors collects field-level validation errors.
type Errors map[string]string

// HasErrors returns true if any validation errors have been recorded.
func (e Errors) HasErrors() bool {
	return len(e) > 0
}

// Required checks that a string field is not empty.
func Required(errors Errors, field, value string) {
	if strings.TrimSpace(value) == "" {
		errors[field] = fmt.Sprintf("%s is required", field)
	}
}

// MinLength checks that a string meets a minimum length.
func MinLength(errors Errors, field, value string, min int) {
	if len(value) < min {
		errors[field] = fmt.Sprintf("%s must be at least %d characters", field, min)
	}
}

// MaxLength checks that a string does not exceed a maximum length.
func MaxLength(errors Errors, field, value string, max int) {
	if len(value) > max {
		errors[field] = fmt.Sprintf("%s must be at most %d characters", field, max)
	}
}

// Email checks that a string is a valid email address.
func Email(errors Errors, field, value string) {
	if value == "" {
		return
	}
	if _, err := mail.ParseAddress(value); err != nil {
		errors[field] = fmt.Sprintf("%s must be a valid email address", field)
	}
}

// MinValue checks that an integer meets a minimum value.
func MinValue(errors Errors, field string, value, min int) {
	if value < min {
		errors[field] = fmt.Sprintf("%s must be at least %d", field, min)
	}
}

// MaxValue checks that an integer does not exceed a maximum value.
func MaxValue(errors Errors, field string, value, max int) {
	if value > max {
		errors[field] = fmt.Sprintf("%s must be at most %d", field, max)
	}
}

// OneOf checks that a string is one of the allowed values.
func OneOf(errors Errors, field, value string, allowed []string) {
	for _, a := range allowed {
		if value == a {
			return
		}
	}
	errors[field] = fmt.Sprintf("%s must be one of: %s", field, strings.Join(allowed, ", "))
}

// PositiveFloat checks that a float is greater than zero.
func PositiveFloat(errors Errors, field string, value float64) {
	if value <= 0 {
		errors[field] = fmt.Sprintf("%s must be a positive number", field)
	}
}

// Password checks that a password meets complexity requirements:
// at least 8 characters, one uppercase letter, one lowercase letter, and one digit.
func Password(errors Errors, field, value string) {
	if len(value) < 8 {
		errors[field] = fmt.Sprintf("%s must be at least 8 characters", field)
		return
	}

	var hasUpper, hasLower, hasDigit bool
	for _, ch := range value {
		switch {
		case unicode.IsUpper(ch):
			hasUpper = true
		case unicode.IsLower(ch):
			hasLower = true
		case unicode.IsDigit(ch):
			hasDigit = true
		}
	}

	var missing []string
	if !hasUpper {
		missing = append(missing, "one uppercase letter")
	}
	if !hasLower {
		missing = append(missing, "one lowercase letter")
	}
	if !hasDigit {
		missing = append(missing, "one number")
	}

	if len(missing) > 0 {
		errors[field] = fmt.Sprintf("%s must contain at least %s", field, strings.Join(missing, ", "))
	}
}
