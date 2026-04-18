package care_logs

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/stride-pro/backend/internal/auth"
	"github.com/stride-pro/backend/internal/subscriptions"
	"github.com/stride-pro/backend/pkg/response"
)

// Handler exposes HTTP endpoints for care log management.
type Handler struct {
	service *Service
}

// NewHandler creates a care log handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// List handles GET /api/horses/{horseId}/care-logs.
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(uuid.UUID)
	horseID, err := uuid.Parse(mux.Vars(r)["horseId"])
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid horse ID")
		return
	}

	logs, err := h.service.GetByHorseID(userID, horseID)
	if err != nil {
		if errors.Is(err, subscriptions.ErrFeatureNotAvailable) {
			response.Error(w, http.StatusForbidden, "Care logs require a Base plan or higher")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Failed to fetch care logs")
		return
	}

	response.JSON(w, http.StatusOK, logs)
}

// Create handles POST /api/horses/{horseId}/care-logs.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(uuid.UUID)
	horseID, err := uuid.Parse(mux.Vars(r)["horseId"])
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid horse ID")
		return
	}

	var input Input
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if errs := input.Validate(); errs.HasErrors() {
		response.ValidationError(w, errs)
		return
	}

	cl, err := h.service.Create(userID, horseID, input)
	if err != nil {
		if errors.Is(err, subscriptions.ErrFeatureNotAvailable) {
			response.Error(w, http.StatusForbidden, "Care logs require a Base plan or higher")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Failed to create care log")
		return
	}

	response.JSON(w, http.StatusCreated, cl)
}

// Update handles PUT /api/care-logs/{id}.
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(uuid.UUID)
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid care log ID")
		return
	}

	var input Input
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if errs := input.Validate(); errs.HasErrors() {
		response.ValidationError(w, errs)
		return
	}

	cl, err := h.service.Update(userID, id, input)
	if err != nil {
		if errors.Is(err, subscriptions.ErrFeatureNotAvailable) {
			response.Error(w, http.StatusForbidden, "Care logs require a Base plan or higher")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Failed to update care log")
		return
	}
	if cl == nil {
		response.Error(w, http.StatusNotFound, "Care log not found")
		return
	}

	response.JSON(w, http.StatusOK, cl)
}

// Delete handles DELETE /api/care-logs/{id}.
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(uuid.UUID)
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid care log ID")
		return
	}

	if err := h.service.Delete(userID, id); err != nil {
		if errors.Is(err, subscriptions.ErrFeatureNotAvailable) {
			response.Error(w, http.StatusForbidden, "Care logs require a Base plan or higher")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Failed to delete care log")
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "Care log deleted"})
}
