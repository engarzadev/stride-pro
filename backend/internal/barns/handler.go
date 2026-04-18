package barns

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

// Handler exposes HTTP endpoints for barn management.
type Handler struct {
	service *Service
}

// NewHandler creates a barn handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// List handles GET /api/barns.
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(uuid.UUID)

	barns, err := h.service.GetAll(userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to fetch barns")
		return
	}
	if barns == nil {
		barns = []models.Barn{}
	}

	response.JSON(w, http.StatusOK, barns)
}

// Get handles GET /api/barns/{id}.
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(uuid.UUID)
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid barn ID")
		return
	}

	barn, err := h.service.GetByID(userID, id)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to fetch barn")
		return
	}
	if barn == nil {
		response.Error(w, http.StatusNotFound, "Barn not found")
		return
	}

	response.JSON(w, http.StatusOK, barn)
}

// Create handles POST /api/barns.
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

	barn, err := h.service.Create(userID, input)
	if err != nil {
		if errors.Is(err, subscriptions.ErrFeatureNotAvailable) {
			response.Error(w, http.StatusForbidden, "Barn management requires a paid plan")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Failed to create barn")
		return
	}

	response.JSON(w, http.StatusCreated, barn)
}

// Update handles PUT /api/barns/{id}.
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(uuid.UUID)
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid barn ID")
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

	barn, err := h.service.Update(userID, id, input)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to update barn")
		return
	}

	response.JSON(w, http.StatusOK, barn)
}

// Delete handles DELETE /api/barns/{id}.
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(uuid.UUID)
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid barn ID")
		return
	}

	if err := h.service.Delete(userID, id); err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to delete barn")
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "Barn deleted"})
}
