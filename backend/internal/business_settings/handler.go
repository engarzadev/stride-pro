package business_settings

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"github.com/stride-pro/backend/internal/auth"
	"github.com/stride-pro/backend/pkg/response"
)

// Handler exposes HTTP endpoints for business settings.
type Handler struct {
	service *Service
}

// NewHandler creates a business settings handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Get handles GET /api/settings/business.
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(auth.UserIDKey).(uuid.UUID)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	bs, err := h.service.Get(userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to fetch business settings")
		return
	}
	response.JSON(w, http.StatusOK, bs)
}

// Upsert handles PUT /api/settings/business.
func (h *Handler) Upsert(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(auth.UserIDKey).(uuid.UUID)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var input UpsertInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	bs, err := h.service.Upsert(userID, input)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to save business settings")
		return
	}
	response.JSON(w, http.StatusOK, bs)
}
