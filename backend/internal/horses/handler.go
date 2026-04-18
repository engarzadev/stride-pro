package horses

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

// Handler exposes HTTP endpoints for horse management.
type Handler struct {
	service *Service
}

// NewHandler creates a horse handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// List handles GET /api/horses.
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(uuid.UUID)

	// Support filtering by client_id or barn_id via query params
	if clientIDStr := r.URL.Query().Get("client_id"); clientIDStr != "" {
		clientID, err := uuid.Parse(clientIDStr)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid client_id")
			return
		}
		horses, err := h.service.GetByClientID(userID, clientID)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "Failed to fetch horses")
			return
		}
		if horses == nil {
			horses = []models.Horse{}
		}
		response.JSON(w, http.StatusOK, horses)
		return
	}

	if barnIDStr := r.URL.Query().Get("barn_id"); barnIDStr != "" {
		barnID, err := uuid.Parse(barnIDStr)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid barn_id")
			return
		}
		horses, err := h.service.GetByBarnID(userID, barnID)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "Failed to fetch horses")
			return
		}
		if horses == nil {
			horses = []models.Horse{}
		}
		response.JSON(w, http.StatusOK, horses)
		return
	}

	horses, err := h.service.GetAll(userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to fetch horses")
		return
	}
	if horses == nil {
		horses = []models.Horse{}
	}

	response.JSON(w, http.StatusOK, horses)
}

// Get handles GET /api/horses/{id}.
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(uuid.UUID)
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid horse ID")
		return
	}

	horse, err := h.service.GetByID(userID, id)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to fetch horse")
		return
	}
	if horse == nil {
		response.Error(w, http.StatusNotFound, "Horse not found")
		return
	}

	response.JSON(w, http.StatusOK, horse)
}

// Create handles POST /api/horses.
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

	horse, err := h.service.Create(userID, input)
	if err != nil {
		if errors.Is(err, subscriptions.ErrLimitExceeded) {
			response.Error(w, http.StatusForbidden, "Horse limit reached for your current plan")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Failed to create horse")
		return
	}

	response.JSON(w, http.StatusCreated, horse)
}

// Update handles PUT /api/horses/{id}.
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(uuid.UUID)
	id, err := uuid.Parse(mux.Vars(r)["id"])
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

	horse, err := h.service.Update(userID, id, input)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to update horse")
		return
	}

	response.JSON(w, http.StatusOK, horse)
}

// Delete handles DELETE /api/horses/{id}.
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(uuid.UUID)
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid horse ID")
		return
	}

	if err := h.service.Delete(userID, id); err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to delete horse")
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "Horse deleted"})
}
