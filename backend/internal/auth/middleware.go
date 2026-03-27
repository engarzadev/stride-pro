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
// the user ID into the request context.
func Middleware(service *Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" {
				response.Error(w, http.StatusUnauthorized, "Authorization header required")
				return
			}

			parts := strings.SplitN(header, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
				response.Error(w, http.StatusUnauthorized, "Invalid authorization header format")
				return
			}

			userID, err := service.ValidateToken(parts[1])
			if err != nil {
				response.Error(w, http.StatusUnauthorized, "Invalid or expired token")
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
