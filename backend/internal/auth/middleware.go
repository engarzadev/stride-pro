package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/stride-pro/backend/pkg/response"
)

type contextKey string

// UserIDKey is the context key used to store the authenticated user's ID.
const UserIDKey contextKey = "userID"

// Middleware returns an HTTP middleware that validates JWT tokens and injects
// the user ID into the request context. It reads the token from the
// `access_token` HttpOnly cookie first, then falls back to the Authorization
// header for API clients that don't use cookies.
func Middleware(service *Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenStr := tokenFromRequest(r)
			if tokenStr == "" {
				response.Error(w, http.StatusUnauthorized, "Authorization required")
				return
			}

			userID, err := service.ValidateToken(tokenStr)
			if err != nil {
				response.Error(w, http.StatusUnauthorized, "Invalid or expired token")
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// tokenFromRequest extracts a JWT from the request, preferring the HttpOnly
// cookie over the Authorization header.
func tokenFromRequest(r *http.Request) string {
	// Prefer cookie — set as HttpOnly so JavaScript cannot read it
	if cookie, err := r.Cookie("access_token"); err == nil && cookie.Value != "" {
		return cookie.Value
	}

	// Fall back to Bearer token for non-browser API clients
	header := r.Header.Get("Authorization")
	if header == "" {
		return ""
	}
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
		return ""
	}
	return parts[1]
}
