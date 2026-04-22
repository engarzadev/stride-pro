package middleware

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"net/http"

	"github.com/stride-pro/backend/pkg/response"
)

const (
	csrfCookieName = "XSRF-TOKEN"
	csrfHeaderName = "X-XSRF-TOKEN"
)

// CSRFSetCookie is applied globally. On every response it ensures the
// XSRF-TOKEN cookie is present so the Angular client can read it.
// The cookie is intentionally NOT HttpOnly — the browser's same-origin
// policy prevents other origins from reading it, and the client needs
// to forward the value as a request header.
func CSRFSetCookie(isProd bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, err := r.Cookie(csrfCookieName); err != nil {
				token, err := generateCSRFToken()
				if err != nil {
					response.Error(w, http.StatusInternalServerError, "Internal server error")
					return
				}
				http.SetCookie(w, &http.Cookie{
					Name:     csrfCookieName,
					Value:    token,
					Path:     "/",
					HttpOnly: false,
					Secure:   isProd,
					SameSite: http.SameSiteLaxMode,
				})
			}

			next.ServeHTTP(w, r)
		})
	}
}

// CSRFValidate is applied to all protected (authenticated) routes. It rejects
// mutating requests that do not supply a valid X-XSRF-TOKEN header matching
// the XSRF-TOKEN cookie. Safe methods (GET, HEAD, OPTIONS) are skipped.
func CSRFValidate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet, http.MethodHead, http.MethodOptions:
			// Safe methods cannot mutate state — no CSRF risk
			next.ServeHTTP(w, r)
			return
		}

		cookie, err := r.Cookie(csrfCookieName)
		if err != nil || cookie.Value == "" {
			response.Error(w, http.StatusForbidden, "CSRF cookie missing")
			return
		}

		header := r.Header.Get(csrfHeaderName)
		if header == "" {
			response.Error(w, http.StatusForbidden, "CSRF token missing")
			return
		}

		// Constant-time comparison prevents timing-based token guessing
		if subtle.ConstantTimeCompare([]byte(header), []byte(cookie.Value)) != 1 {
			response.Error(w, http.StatusForbidden, "CSRF token invalid")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func generateCSRFToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
