package middleware

import (
	"net/http"
)

// HTTPSRedirect returns middleware that enforces HTTPS when the application runs
// behind a TLS-terminating proxy (e.g. AWS ALB, Cloudflare, nginx). The proxy
// is expected to set X-Forwarded-Proto on every request. Any request where that
// header is not "https" receives a permanent 301 redirect to the HTTPS equivalent.
func HTTPSRedirect(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Forwarded-Proto") != "https" {
			target := "https://" + r.Host + r.URL.RequestURI()
			http.Redirect(w, r, target, http.StatusMovedPermanently)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// SecurityHeaders returns middleware that sets security-related HTTP response headers
// on every request. Pass isProd=true to also set Strict-Transport-Security (HSTS).
func SecurityHeaders(isProd bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Prevent clickjacking
			w.Header().Set("X-Frame-Options", "DENY")
			// Prevent MIME-type sniffing
			w.Header().Set("X-Content-Type-Options", "nosniff")
			// Legacy XSS filter (belt-and-suspenders for older browsers)
			w.Header().Set("X-XSS-Protection", "1; mode=block")
			// Limit referrer information sent to third parties
			w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
			// Only send HSTS in production — it cannot easily be undone once set
			if isProd {
				w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
			}

			next.ServeHTTP(w, r)
		})
	}
}
