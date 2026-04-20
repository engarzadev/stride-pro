package service_items

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/stride-pro/backend/internal/auth"
	"github.com/stride-pro/backend/internal/models"
	"github.com/stride-pro/backend/pkg/response"
)

// Handler exposes HTTP endpoints for service item management.
type Handler struct {
	service *Service
}

// NewHandler creates a service item handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// List handles GET /api/settings/service-items.
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(auth.UserIDKey).(uuid.UUID)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	items, err := h.service.GetAll(userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to fetch service items")
		return
	}
	if items == nil {
		items = []models.ServiceItem{}
	}
	response.JSON(w, http.StatusOK, items)
}

// Create handles POST /api/settings/service-items.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(auth.UserIDKey).(uuid.UUID)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var input ItemInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if input.Name == "" {
		response.Error(w, http.StatusBadRequest, "Name is required")
		return
	}

	item, err := h.service.Create(userID, input)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to create service item")
		return
	}
	response.JSON(w, http.StatusCreated, item)
}

// Update handles PUT /api/settings/service-items/{id}.
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(auth.UserIDKey).(uuid.UUID)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid item ID")
		return
	}

	var input ItemInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if input.Name == "" {
		response.Error(w, http.StatusBadRequest, "Name is required")
		return
	}

	item, err := h.service.Update(userID, id, input)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to update service item")
		return
	}
	response.JSON(w, http.StatusOK, item)
}

// Delete handles DELETE /api/settings/service-items/{id}.
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(auth.UserIDKey).(uuid.UUID)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid item ID")
		return
	}

	if err := h.service.Delete(userID, id); err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to delete service item")
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "Service item deleted"})
}
