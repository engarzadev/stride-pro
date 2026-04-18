package subscriptions

import (
	"net/http"

	"github.com/google/uuid"

	"github.com/stride-pro/backend/internal/auth"
	"github.com/stride-pro/backend/pkg/response"
)

// Handler exposes HTTP endpoints for subscription information.
type Handler struct {
	service *Service
}

// NewHandler creates a subscription handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

type usageCount struct {
	Count int `json:"count"`
	Limit int `json:"limit"` // -1 means unlimited
}

type subscriptionResponse struct {
	Plan  *Plan `json:"plan"`
	Usage struct {
		Clients usageCount `json:"clients"`
		Horses  usageCount `json:"horses"`
	} `json:"usage"`
}

// Get handles GET /api/subscription and returns the user's current plan and usage.
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(uuid.UUID)

	plan, err := h.service.GetCurrentPlan(userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to fetch subscription")
		return
	}

	clientLimit, err := h.service.GetClientLimit(userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to fetch subscription")
		return
	}

	horseLimit, err := h.service.GetHorseLimit(userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to fetch subscription")
		return
	}

	var clientCount, horseCount int
	h.service.db.QueryRow("SELECT COUNT(*) FROM clients WHERE user_id = $1", userID).Scan(&clientCount)
	h.service.db.QueryRow("SELECT COUNT(*) FROM horses WHERE user_id = $1", userID).Scan(&horseCount)

	resp := subscriptionResponse{Plan: plan}
	resp.Usage.Clients = usageCount{Count: clientCount, Limit: clientLimit}
	resp.Usage.Horses = usageCount{Count: horseCount, Limit: horseLimit}

	response.JSON(w, http.StatusOK, resp)
}
