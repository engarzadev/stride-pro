package reminders

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/stride-pro/backend/internal/auth"
	"github.com/stride-pro/backend/pkg/response"
)

// Handler exposes HTTP endpoints for reminder management.
type Handler struct {
	service *Service
}

// NewHandler creates a reminder handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// List handles GET /api/horses/{horseId}/reminders.
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(auth.UserIDKey).(uuid.UUID)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	horseID, err := uuid.Parse(mux.Vars(r)["horseId"])
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid horse ID")
		return
	}

	list, err := h.service.GetByHorseID(userID, horseID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to fetch reminders")
		return
	}

	response.JSON(w, http.StatusOK, list)
}

// Create handles POST /api/horses/{horseId}/reminders.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(auth.UserIDKey).(uuid.UUID)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	horseID, err := uuid.Parse(mux.Vars(r)["horseId"])
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid horse ID")
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

	rm, err := h.service.Create(userID, horseID, input)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to create reminder")
		return
	}

	response.JSON(w, http.StatusCreated, rm)
}

// Update handles PUT /api/reminders/{id}.
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(auth.UserIDKey).(uuid.UUID)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid reminder ID")
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

	rm, err := h.service.Update(userID, id, input)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to update reminder")
		return
	}
	if rm == nil {
		response.Error(w, http.StatusNotFound, "Reminder not found")
		return
	}

	response.JSON(w, http.StatusOK, rm)
}

// Patch handles PATCH /api/reminders/{id}.
func (h *Handler) Patch(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(auth.UserIDKey).(uuid.UUID)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid reminder ID")
		return
	}

	var input PatchInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	rm, err := h.service.Patch(userID, id, input)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to update reminder")
		return
	}
	if rm == nil {
		response.Error(w, http.StatusNotFound, "Reminder not found")
		return
	}

	response.JSON(w, http.StatusOK, rm)
}

// Delete handles DELETE /api/reminders/{id}.
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(auth.UserIDKey).(uuid.UUID)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid reminder ID")
		return
	}

	if err := h.service.Delete(userID, id); err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to delete reminder")
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "Reminder deleted"})
}
