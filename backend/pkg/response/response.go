// Package response provides helpers for consistent JSON API responses.
package response

import (
	"encoding/json"
	"net/http"
)

// Envelope is the standard API response wrapper.
type Envelope struct {
	Data  interface{} `json:"data,omitempty"`
	Error *APIError   `json:"error,omitempty"`
	Meta  *Meta       `json:"meta,omitempty"`
}

// APIError describes an error returned to the client.
type APIError struct {
	Message string            `json:"message"`
	Code    string            `json:"code,omitempty"`
	Fields  map[string]string `json:"fields,omitempty"`
}

// Meta carries optional pagination or request metadata.
type Meta struct {
	Total  int `json:"total,omitempty"`
	Page   int `json:"page,omitempty"`
	Limit  int `json:"limit,omitempty"`
}

// JSON writes a success response with the given status code.
func JSON(w http.ResponseWriter, status int, data interface{}) {
	writeJSON(w, status, Envelope{Data: data})
}

// JSONWithMeta writes a success response with metadata.
func JSONWithMeta(w http.ResponseWriter, status int, data interface{}, meta *Meta) {
	writeJSON(w, status, Envelope{Data: data, Meta: meta})
}

// Error writes an error response with the given status code and message.
func Error(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, Envelope{
		Error: &APIError{Message: message},
	})
}

// ErrorWithCode writes an error response with a machine-readable code.
func ErrorWithCode(w http.ResponseWriter, status int, message, code string) {
	writeJSON(w, status, Envelope{
		Error: &APIError{Message: message, Code: code},
	})
}

// ValidationError writes a 422 response with per-field validation errors.
func ValidationError(w http.ResponseWriter, fields map[string]string) {
	writeJSON(w, http.StatusUnprocessableEntity, Envelope{
		Error: &APIError{
			Message: "Validation failed",
			Code:    "VALIDATION_ERROR",
			Fields:  fields,
		},
	})
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
