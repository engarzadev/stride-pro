// Package router configures all HTTP routes for the application.
package router

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"

	"github.com/stride-pro/backend/internal/appointments"
	"github.com/stride-pro/backend/internal/auth"
	"github.com/stride-pro/backend/internal/barns"
	biz "github.com/stride-pro/backend/internal/business_settings"
	carelogs "github.com/stride-pro/backend/internal/care_logs"
	"github.com/stride-pro/backend/internal/clients"
	"github.com/stride-pro/backend/internal/config"
	"github.com/stride-pro/backend/internal/database"
	"github.com/stride-pro/backend/internal/horses"
	"github.com/stride-pro/backend/internal/invoices"
	"github.com/stride-pro/backend/internal/middleware"
	"github.com/stride-pro/backend/internal/reminders"
	"github.com/stride-pro/backend/internal/sessions"
	svc "github.com/stride-pro/backend/internal/service_items"
	"github.com/stride-pro/backend/internal/subscriptions"
	"github.com/stride-pro/backend/pkg/response"
)

// Deps holds all handler dependencies needed to configure routing.
type Deps struct {
	DB                     *database.DB
	Config                 *config.Config
	AuthService            *auth.Service
	AuthHandler            *auth.Handler
	ClientHandler          *clients.Handler
	HorseHandler           *horses.Handler
	BarnHandler            *barns.Handler
	ApptHandler            *appointments.Handler
	SessionHandler         *sessions.Handler
	InvoiceHandler         *invoices.Handler
	SubscriptionHandler    *subscriptions.Handler
	BusinessSettingHandler *biz.Handler
	ServiceItemHandler     *svc.Handler
	CareLogHandler         *carelogs.Handler
	ReminderHandler        *reminders.Handler
}

// New creates and configures the application router with all routes and middleware.
func New(deps Deps) http.Handler {
	r := mux.NewRouter()

	// Recovery must be outermost — catches panics from all middleware and handlers below
	r.Use(middleware.Recovery)

	// CORS must run before HTTPSRedirect so that redirect responses (301) also carry
	// Access-Control-Allow-Origin. Without this, the browser blocks the redirect and
	// the request fails with a CORS error before it can reach the actual handler.
	// X-XSRF-TOKEN is included so the Angular CSRF header is not blocked by preflight.
	corsConfig := middleware.CORSConfig{
		AllowedOrigins: deps.Config.AllowedOrigins,
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Authorization", "Content-Type", "Accept", "X-XSRF-TOKEN"},
	}
	r.Use(middleware.CORS(corsConfig))

	// In proxy mode, redirect plain HTTP after CORS headers are already set
	if deps.Config.TLSProxyMode {
		r.Use(middleware.HTTPSRedirect)
	}

	// Security headers on every response
	r.Use(middleware.SecurityHeaders(deps.Config.IsProd()))

	// Set XSRF-TOKEN cookie on every response so the Angular client always has it
	r.Use(middleware.CSRFSetCookie(deps.Config.IsProd()))

	r.Use(middleware.Logging)

	// Global rate limiter: 10 req/s per IP, burst of 50
	globalLimiter := middleware.NewRateLimiter(10, 50)
	r.Use(globalLimiter.Middleware)

	// Request body size limit: 2 MB
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, 2<<20) // 2 MB
			next.ServeHTTP(w, r)
		})
	})

	// Handle preflight OPTIONS requests for all routes
	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}).Methods("OPTIONS")

	// Health check
	r.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		if err := deps.DB.HealthCheck(); err != nil {
			response.Error(w, http.StatusServiceUnavailable, "Database unavailable")
			return
		}
		response.JSON(w, http.StatusOK, map[string]string{"status": "healthy"})
	}).Methods("GET")

	// Auth routes — stricter rate limit (5 req/min per IP)
	authLimiter := middleware.NewAuthRateLimiter()
	r.Handle("/api/auth/register", authLimiter.Middleware(http.HandlerFunc(deps.AuthHandler.Register))).Methods("POST")
	r.Handle("/api/auth/login", authLimiter.Middleware(http.HandlerFunc(deps.AuthHandler.Login))).Methods("POST")
	r.HandleFunc("/api/auth/refresh", deps.AuthHandler.Refresh).Methods("POST")

	// Protected routes
	protected := r.PathPrefix("/api").Subrouter()
	protected.Use(auth.Middleware(deps.AuthService))
	// CSRF validation on all mutating requests within protected routes
	protected.Use(middleware.CSRFValidate)

	// Auth - protected
	protected.HandleFunc("/auth/me", deps.AuthHandler.Me).Methods("GET")
	protected.HandleFunc("/auth/logout", deps.AuthHandler.Logout).Methods("POST")

	// Subscription
	protected.HandleFunc("/subscription", deps.SubscriptionHandler.Get).Methods("GET")

	// Clients
	protected.HandleFunc("/clients", deps.ClientHandler.List).Methods("GET")
	protected.HandleFunc("/clients/{id}", deps.ClientHandler.Get).Methods("GET")
	protected.HandleFunc("/clients", deps.ClientHandler.Create).Methods("POST")
	protected.HandleFunc("/clients/{id}", deps.ClientHandler.Update).Methods("PUT")
	protected.HandleFunc("/clients/{id}", deps.ClientHandler.Delete).Methods("DELETE")

	// Horses
	protected.HandleFunc("/horses", deps.HorseHandler.List).Methods("GET")
	protected.HandleFunc("/horses/{id}", deps.HorseHandler.Get).Methods("GET")
	protected.HandleFunc("/horses", deps.HorseHandler.Create).Methods("POST")
	protected.HandleFunc("/horses/{id}", deps.HorseHandler.Update).Methods("PUT")
	protected.HandleFunc("/horses/{id}", deps.HorseHandler.Delete).Methods("DELETE")

	// Barns
	protected.HandleFunc("/barns", deps.BarnHandler.List).Methods("GET")
	protected.HandleFunc("/barns/{id}", deps.BarnHandler.Get).Methods("GET")
	protected.HandleFunc("/barns", deps.BarnHandler.Create).Methods("POST")
	protected.HandleFunc("/barns/{id}", deps.BarnHandler.Update).Methods("PUT")
	protected.HandleFunc("/barns/{id}", deps.BarnHandler.Delete).Methods("DELETE")

	// Appointments
	protected.HandleFunc("/appointments", deps.ApptHandler.List).Methods("GET")
	protected.HandleFunc("/appointments/{id}", deps.ApptHandler.Get).Methods("GET")
	protected.HandleFunc("/appointments", deps.ApptHandler.Create).Methods("POST")
	protected.HandleFunc("/appointments/{id}", deps.ApptHandler.Update).Methods("PUT")
	protected.HandleFunc("/appointments/{id}", deps.ApptHandler.Delete).Methods("DELETE")

	// Sessions
	protected.HandleFunc("/sessions", deps.SessionHandler.List).Methods("GET")
	protected.HandleFunc("/sessions/{id}", deps.SessionHandler.Get).Methods("GET")
	protected.HandleFunc("/sessions", deps.SessionHandler.Create).Methods("POST")
	protected.HandleFunc("/sessions/{id}", deps.SessionHandler.Update).Methods("PUT")
	protected.HandleFunc("/sessions/{id}", deps.SessionHandler.Delete).Methods("DELETE")

	// Invoices
	protected.HandleFunc("/invoices", deps.InvoiceHandler.List).Methods("GET")
	protected.HandleFunc("/invoices/{id}", deps.InvoiceHandler.Get).Methods("GET")
	protected.HandleFunc("/invoices", deps.InvoiceHandler.Create).Methods("POST")
	protected.HandleFunc("/invoices/{id}", deps.InvoiceHandler.Update).Methods("PUT")
	protected.HandleFunc("/invoices/{id}", deps.InvoiceHandler.Delete).Methods("DELETE")
	protected.HandleFunc("/invoices/{id}/status", deps.InvoiceHandler.UpdateStatus).Methods("PATCH")
	protected.HandleFunc("/invoices/{id}/send", deps.InvoiceHandler.Send).Methods("POST")

	// Care logs
	protected.HandleFunc("/horses/{horseId}/care-logs", deps.CareLogHandler.List).Methods("GET")
	protected.HandleFunc("/horses/{horseId}/care-logs", deps.CareLogHandler.Create).Methods("POST")
	protected.HandleFunc("/care-logs/{id}", deps.CareLogHandler.Update).Methods("PUT")
	protected.HandleFunc("/care-logs/{id}", deps.CareLogHandler.Delete).Methods("DELETE")

	// Reminders
	protected.HandleFunc("/horses/{horseId}/reminders", deps.ReminderHandler.List).Methods("GET")
	protected.HandleFunc("/horses/{horseId}/reminders", deps.ReminderHandler.Create).Methods("POST")
	protected.HandleFunc("/reminders/{id}", deps.ReminderHandler.Update).Methods("PUT")
	protected.HandleFunc("/reminders/{id}", deps.ReminderHandler.Patch).Methods("PATCH")
	protected.HandleFunc("/reminders/{id}", deps.ReminderHandler.Delete).Methods("DELETE")

	// Settings
	protected.HandleFunc("/settings/business", deps.BusinessSettingHandler.Get).Methods("GET")
	protected.HandleFunc("/settings/business", deps.BusinessSettingHandler.Upsert).Methods("PUT")
	protected.HandleFunc("/settings/service-items", deps.ServiceItemHandler.List).Methods("GET")
	protected.HandleFunc("/settings/service-items", deps.ServiceItemHandler.Create).Methods("POST")
	protected.HandleFunc("/settings/service-items/{id}", deps.ServiceItemHandler.Update).Methods("PUT")
	protected.HandleFunc("/settings/service-items/{id}", deps.ServiceItemHandler.Delete).Methods("DELETE")

	// Serve the Angular SPA in production — static files from the embedded
	// dist directory, with a fallback to index.html for client-side routes.
	spaDir := filepath.Join(".", "public")
	if _, err := os.Stat(spaDir); err == nil {
		spa := spaHandler{staticPath: spaDir, indexPath: "index.html"}
		r.PathPrefix("/").Handler(spa)
	}

	return r
}

// spaHandler serves static files from a directory and falls back to index.html
// for any path that doesn't match a file, enabling Angular's client-side routing.
type spaHandler struct {
	staticPath string
	indexPath  string
}

func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join(h.staticPath, filepath.Clean(r.URL.Path))

	if info, err := os.Stat(path); err == nil && !info.IsDir() {
		http.ServeFile(w, r, path)
		return
	}

	http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
}
