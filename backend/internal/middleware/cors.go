package middleware

import (
	"net/http"
	"strings"
)

// CORSConfig holds CORS middleware configuration.
type CORSConfig struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
}

// DefaultCORSConfig returns a development CORS configuration allowing localhost only.
func DefaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowedOrigins: []string{"http://localhost:4200"},
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Authorization", "Content-Type", "Accept"},
	}
}

// CORS returns middleware that validates the request Origin against the allowed list
// and reflects it back only when it matches. This prevents cross-origin requests from
// unknown domains and stops CSRF attacks that relied on the old wildcard behaviour.
func CORS(cfg CORSConfig) func(http.Handler) http.Handler {
	methods := strings.Join(cfg.AllowedMethods, ", ")
	headers := strings.Join(cfg.AllowedHeaders, ", ")

	// Build a lookup set for O(1) origin checks
	originSet := make(map[string]struct{}, len(cfg.AllowedOrigins))
	for _, o := range cfg.AllowedOrigins {
		originSet[strings.ToLower(o)] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			if origin != "" {
				if _, ok := originSet[strings.ToLower(origin)]; ok {
					// Reflect the exact origin back (never a wildcard)
					w.Header().Set("Access-Control-Allow-Origin", origin)
					w.Header().Set("Access-Control-Allow-Methods", methods)
					w.Header().Set("Access-Control-Allow-Headers", headers)
					w.Header().Set("Access-Control-Max-Age", "86400")
					// Required for HttpOnly cookies to be sent cross-origin
					w.Header().Set("Access-Control-Allow-Credentials", "true")
					// Vary must be set so caches don't serve the wrong origin
					w.Header().Add("Vary", "Origin")
				}
				// Unrecognised origin — headers are simply not set; the browser
				// will block the response automatically.
			}

			// Handle preflight requests
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
