package invoices

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/stride-pro/backend/internal/auth"
	"github.com/stride-pro/backend/internal/models"
	"github.com/stride-pro/backend/pkg/response"
	"github.com/stride-pro/backend/pkg/validator"
)

// Handler exposes HTTP endpoints for invoice management.
type Handler struct {
	service *Service
}

// NewHandler creates an invoice handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// List handles GET /api/invoices. Supports optional client_id query param.
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(uuid.UUID)

	if clientIDStr := r.URL.Query().Get("client_id"); clientIDStr != "" {
		clientID, err := uuid.Parse(clientIDStr)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid client_id")
			return
		}
		invoices, err := h.service.GetByClientID(userID, clientID)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "Failed to fetch invoices")
			return
		}
		if invoices == nil {
			invoices = []models.Invoice{}
		}
		response.JSON(w, http.StatusOK, invoices)
		return
	}

	invoices, err := h.service.GetAll(userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to fetch invoices")
		return
	}
	if invoices == nil {
		invoices = []models.Invoice{}
	}

	response.JSON(w, http.StatusOK, invoices)
}

// Get handles GET /api/invoices/{id}.
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(uuid.UUID)
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid invoice ID")
		return
	}

	invoice, err := h.service.GetByID(userID, id)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to fetch invoice")
		return
	}
	if invoice == nil {
		response.Error(w, http.StatusNotFound, "Invoice not found")
		return
	}

	response.JSON(w, http.StatusOK, invoice)
}

// Create handles POST /api/invoices.
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

	invoice, err := h.service.Create(userID, input)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to create invoice")
		return
	}

	response.JSON(w, http.StatusCreated, invoice)
}

// Update handles PUT /api/invoices/{id}.
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(uuid.UUID)
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid invoice ID")
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

	invoice, err := h.service.Update(userID, id, input)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to update invoice")
		return
	}

	response.JSON(w, http.StatusOK, invoice)
}

// UpdateStatus handles PATCH /api/invoices/{id}/status.
func (h *Handler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(uuid.UUID)
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid invoice ID")
		return
	}

	var body struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	errs := validator.Errors{}
	validator.OneOf(errs, "status", body.Status, []string{"draft", "sent", "paid", "overdue", "cancelled"})
	if errs.HasErrors() {
		response.ValidationError(w, errs)
		return
	}

	if err := h.service.UpdateStatus(userID, id, body.Status); err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to update invoice status")
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "Invoice status updated"})
}

// Send handles POST /api/invoices/{id}/send — emails the invoice and marks it sent.
func (h *Handler) Send(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(uuid.UUID)
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid invoice ID")
		return
	}

	if err := h.service.SendInvoice(userID, id); err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to send invoice")
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "Invoice sent"})
}

// Delete handles DELETE /api/invoices/{id}.
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(uuid.UUID)
	id, err := uuid.Parse(mux.Vars(r)["id"])
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid invoice ID")
		return
	}

	if err := h.service.Delete(userID, id); err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to delete invoice")
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "Invoice deleted"})
}
