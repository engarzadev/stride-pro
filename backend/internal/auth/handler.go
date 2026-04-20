package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/stride-pro/backend/pkg/response"
	"github.com/stride-pro/backend/pkg/validator"
)

// Handler exposes HTTP endpoints for authentication.
type Handler struct {
	service *Service
	isProd  bool
}

// NewHandler creates an auth handler.
func NewHandler(service *Service, isProd bool) *Handler {
	return &Handler{service: service, isProd: isProd}
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
	validator.Password(errs, "password", input.Password)
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

	h.setTokenCookies(w, tokens)

	// Only return the user and expiry — token strings stay in HttpOnly cookies
	response.JSON(w, http.StatusCreated, map[string]interface{}{
		"user":       user,
		"expires_at": tokens.ExpiresAt,
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

	h.setTokenCookies(w, tokens)

	// Only return the user and expiry — token strings stay in HttpOnly cookies
	response.JSON(w, http.StatusOK, map[string]interface{}{
		"user":       user,
		"expires_at": tokens.ExpiresAt,
	})
}

// Logout handles POST /api/auth/logout. Revokes the current access token
// and clears the auth cookies.
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	// Try to revoke the token from cookie first, then fall back to header
	tokenStr := h.tokenFromCookie(r)
	if tokenStr == "" {
		header := r.Header.Get("Authorization")
		parts := strings.SplitN(header, " ", 2)
		if len(parts) == 2 && strings.EqualFold(parts[0], "bearer") {
			tokenStr = parts[1]
		}
	}

	if tokenStr != "" {
		// Best-effort revocation — don't fail the logout if this errors
		_ = h.service.RevokeToken(tokenStr)
	}

	h.clearTokenCookies(w)
	response.JSON(w, http.StatusOK, map[string]string{"message": "Logged out successfully"})
}

// Refresh handles POST /api/auth/refresh.
func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	// Accept refresh token from cookie or request body
	refreshToken := h.refreshTokenFromCookie(r)

	if refreshToken == "" {
		var body struct {
			RefreshToken string `json:"refresh_token"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err == nil {
			refreshToken = body.RefreshToken
		}
	}

	if refreshToken == "" {
		response.Error(w, http.StatusBadRequest, "refresh_token is required")
		return
	}

	tokens, err := h.service.RefreshToken(refreshToken)
	if err != nil {
		h.clearTokenCookies(w)
		response.Error(w, http.StatusUnauthorized, "Invalid or expired refresh token")
		return
	}

	h.setTokenCookies(w, tokens)
	response.JSON(w, http.StatusOK, map[string]interface{}{
		"expires_at": tokens.ExpiresAt,
	})
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

// cookieSameSite returns the appropriate SameSite policy for the environment.
// In production the frontend and backend are on different domains (Vercel vs Railway),
// so cross-site requests require SameSite=None + Secure=true.
// In development both run on localhost so Lax is sufficient and avoids the
// browser requirement for Secure on None cookies.
func (h *Handler) cookieSameSite() http.SameSite {
	if h.isProd {
		return http.SameSiteNoneMode
	}
	return http.SameSiteLaxMode
}

// setTokenCookies writes access and refresh tokens as HttpOnly cookies.
func (h *Handler) setTokenCookies(w http.ResponseWriter, tokens *TokenPair) {
	accessExpiry := time.Unix(tokens.ExpiresAt, 0)
	refreshExpiry := time.Now().Add(7 * 24 * time.Hour)
	sameSite := h.cookieSameSite()

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    tokens.AccessToken,
		Path:     "/",
		Expires:  accessExpiry,
		HttpOnly: true,
		Secure:   h.isProd,
		SameSite: sameSite,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		Path:     "/api/auth/refresh",
		Expires:  refreshExpiry,
		HttpOnly: true,
		Secure:   h.isProd,
		SameSite: sameSite,
	})
}

// clearTokenCookies expires both auth cookies immediately.
func (h *Handler) clearTokenCookies(w http.ResponseWriter) {
	sameSite := h.cookieSameSite()

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   h.isProd,
		SameSite: sameSite,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/api/auth/refresh",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   h.isProd,
		SameSite: sameSite,
	})
}

func (h *Handler) tokenFromCookie(r *http.Request) string {
	c, err := r.Cookie("access_token")
	if err != nil {
		return ""
	}
	return c.Value
}

func (h *Handler) refreshTokenFromCookie(r *http.Request) string {
	c, err := r.Cookie("refresh_token")
	if err != nil {
		return ""
	}
	return c.Value
}
