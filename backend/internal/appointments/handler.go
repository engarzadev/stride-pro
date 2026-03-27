package appointments

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/stride-pro/backend/internal/auth"
	"github.com/stride-pro/backend/internal/models"
	"github.com/stride-pro/backend/pkg/response"
)

// Handler exposes HTTP endpoints for appointment management.
type Handler struct {
	service *Service
}

// NewHandler creates an appointment handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// List handles GET /api/appointments. Supports optional start/end query params for date filtering.
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(uuid.UUID)

	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")

	if startStr != "" && endStr != "" {
		start, err := time.Parse(time.RFC3339, startStr)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid start date format, use RFC3339")
			return
		}
		end, err := time.Parse(time.RFC3339, endStr)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid end date format, use RFC3339")
			return
		}
		appts, err := h.service.GetByDateRange(userID, start, end)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "Failed to fetch appointments")
			return
		}
		if appts == nil {
			appts = []models.Appointment{}
		}
		response.JSON(w, http.StatusOK, appts)
		return
	}

	appts, err := h.service.GetAll(userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to fetch appointments")
		return
	}
	if appts == nil {
		appts = []models.Appointment{}
	}

	response.JSON(w, http.StatusOK, appts)
}

// Get handles GET /api/appointments/{id}.
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(uuid.UUID)
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid appointment ID")
		return
	}

	appt, err := h.service.GetByID(userID, id)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to fetch appointment")
		return
	}
	if appt == nil {
		response.Error(w, http.StatusNotFound, "Appointment not found")
		return
	}

	response.JSON(w, http.StatusOK, appt)
}

// Create handles POST /api/appointments.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(uuid.UUID)

	var input CreateInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if errs := input.Validate(); errs.HasErrors() {
		response.ValidationError(w, errs)
		return
	}

	appt, err := h.service.Create(userID, input)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to create appointment")
		return
	}

	response.JSON(w, http.StatusCreated, appt)
}

// Update handles PUT /api/appointments/{id}.
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(uuid.UUID)
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid appointment ID")
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

	appt, err := h.service.Update(userID, id, input)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to update appointment")
		return
	}

	response.JSON(w, http.StatusOK, appt)
}

// Delete handles DELETE /api/appointments/{id}.
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(uuid.UUID)
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid appointment ID")
		return
	}

	if err := h.service.Delete(userID, id); err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to delete appointment")
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "Appointment deleted"})
}
