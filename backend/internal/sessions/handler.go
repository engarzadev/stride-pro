package sessions

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/stride-pro/backend/internal/auth"
	"github.com/stride-pro/backend/internal/models"
	"github.com/stride-pro/backend/internal/subscriptions"
	"github.com/stride-pro/backend/pkg/response"
)

// Handler exposes HTTP endpoints for session management.
type Handler struct {
	service *Service
}

// NewHandler creates a session handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// List handles GET /api/sessions. Supports optional appointment_id query param.
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(auth.UserIDKey).(uuid.UUID)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if apptIDStr := r.URL.Query().Get("appointment_id"); apptIDStr != "" {
		apptID, err := uuid.Parse(apptIDStr)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid appointment_id")
			return
		}
		sessions, err := h.service.GetByAppointmentID(userID, apptID)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "Failed to fetch sessions")
			return
		}
		if sessions == nil {
			sessions = []models.Session{}
		}
		response.JSON(w, http.StatusOK, sessions)
		return
	}

	sessions, err := h.service.GetAll(userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to fetch sessions")
		return
	}
	if sessions == nil {
		sessions = []models.Session{}
	}

	response.JSON(w, http.StatusOK, sessions)
}

// Get handles GET /api/sessions/{id}.
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(auth.UserIDKey).(uuid.UUID)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid session ID")
		return
	}

	session, err := h.service.GetByID(userID, id)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to fetch session")
		return
	}
	if session == nil {
		response.Error(w, http.StatusNotFound, "Session not found")
		return
	}

	response.JSON(w, http.StatusOK, session)
}

// Create handles POST /api/sessions.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(auth.UserIDKey).(uuid.UUID)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var input CreateInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if errs := input.Validate(); errs.HasErrors() {
		response.ValidationError(w, errs)
		return
	}

	session, err := h.service.Create(userID, input)
	if err != nil {
		if errors.Is(err, subscriptions.ErrFeatureNotAvailable) {
			response.Error(w, http.StatusForbidden, "Session notes require a paid plan")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Failed to create session")
		return
	}

	response.JSON(w, http.StatusCreated, session)
}

// Update handles PUT /api/sessions/{id}.
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(auth.UserIDKey).(uuid.UUID)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid session ID")
		return
	}

	var input CreateInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if errs := input.Validate(); errs.HasErrors() {
		response.ValidationError(w, errs)
		return
	}

	session, err := h.service.Update(userID, id, input)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to update session")
		return
	}

	response.JSON(w, http.StatusOK, session)
}

// Delete handles DELETE /api/sessions/{id}.
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(auth.UserIDKey).(uuid.UUID)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid session ID")
		return
	}

	if err := h.service.Delete(userID, id); err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to delete session")
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "Session deleted"})
}
