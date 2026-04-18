package auth

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"

	"github.com/stride-pro/backend/pkg/response"
	"github.com/stride-pro/backend/pkg/validator"
)

// Handler exposes HTTP endpoints for authentication.
type Handler struct {
	service *Service
}

// NewHandler creates an auth handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Register handles POST /api/auth/register.
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var input RegisterInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	errs := validator.Errors{}
	validator.Required(errs, "email", input.Email)
	validator.Email(errs, "email", input.Email)
	validator.Required(errs, "password", input.Password)
	validator.MinLength(errs, "password", input.Password, 8)
	validator.Required(errs, "first_name", input.FirstName)
	validator.Required(errs, "last_name", input.LastName)
	if errs.HasErrors() {
		response.ValidationError(w, errs)
		return
	}

	user, tokens, err := h.service.Register(input)
	if errors.Is(err, ErrEmailTaken) {
		response.ErrorWithCode(w, http.StatusConflict, "Email already registered", "EMAIL_TAKEN")
		return
	}
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to register user")
		return
	}

	response.JSON(w, http.StatusCreated, map[string]interface{}{
		"user":   user,
		"tokens": tokens,
	})
}

// Login handles POST /api/auth/login.
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var input LoginInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	errs := validator.Errors{}
	validator.Required(errs, "email", input.Email)
	validator.Required(errs, "password", input.Password)
	if errs.HasErrors() {
		response.ValidationError(w, errs)
		return
	}

	user, tokens, err := h.service.Login(input)
	if errors.Is(err, ErrInvalidCredentials) {
		response.Error(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Failed to log in")
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"user":   user,
		"tokens": tokens,
	})
}

// Refresh handles POST /api/auth/refresh.
func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if body.RefreshToken == "" {
		response.Error(w, http.StatusBadRequest, "refresh_token is required")
		return
	}

	tokens, err := h.service.RefreshToken(body.RefreshToken)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "Invalid or expired refresh token")
		return
	}

	response.JSON(w, http.StatusOK, tokens)
}

// Me handles GET /api/auth/me and returns the authenticated user.
func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(UserIDKey).(uuid.UUID)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	user, err := h.service.GetUserByID(userID)
	if err != nil {
		response.Error(w, http.StatusNotFound, "User not found")
		return
	}

	response.JSON(w, http.StatusOK, user)
}
