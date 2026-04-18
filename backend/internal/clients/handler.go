package clients

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

// Handler exposes HTTP endpoints for client management.
type Handler struct {
	service *Service
}

// NewHandler creates a client handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// List handles GET /api/clients.
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(uuid.UUID)

	clients, err := h.service.GetAll(userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to fetch clients")
		return
	}
	if clients == nil {
		clients = []models.Client{}
	}

	response.JSON(w, http.StatusOK, clients)
}

// Get handles GET /api/clients/{id}.
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(uuid.UUID)
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid client ID")
		return
	}

	client, err := h.service.GetByID(userID, id)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to fetch client")
		return
	}
	if client == nil {
		response.Error(w, http.StatusNotFound, "Client not found")
		return
	}

	response.JSON(w, http.StatusOK, client)
}

// Create handles POST /api/clients.
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

	client, err := h.service.Create(userID, input)
	if err != nil {
		if errors.Is(err, subscriptions.ErrLimitExceeded) {
			response.Error(w, http.StatusForbidden, "Client limit reached for your current plan")
			return
		}
		response.Error(w, http.StatusInternalServerError, "Failed to create client")
		return
	}

	response.JSON(w, http.StatusCreated, client)
}

// Update handles PUT /api/clients/{id}.
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(uuid.UUID)
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid client ID")
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

	client, err := h.service.Update(userID, id, input)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to update client")
		return
	}

	response.JSON(w, http.StatusOK, client)
}

// Delete handles DELETE /api/clients/{id}.
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(uuid.UUID)
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid client ID")
		return
	}

	if err := h.service.Delete(userID, id); err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to delete client")
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "Client deleted"})
}
