package middleware

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/stride-pro/backend/pkg/response"
)

// Recovery returns middleware that catches any panic in a downstream handler,
// logs the stack trace server-side, and returns a generic 500 to the client.
// It must be the outermost middleware so it covers all handlers and middleware
// below it in the chain.
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Log full stack trace for debugging — never sent to the client
				log.Printf("PANIC recovered: %v\n%s", err, debug.Stack())
				response.Error(w, http.StatusInternalServerError, "Internal server error")
			}
		}()

		next.ServeHTTP(w, r)
	})
}
