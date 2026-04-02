// Package router configures all HTTP routes for the application.
package router

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/stride-pro/backend/internal/appointments"
	"github.com/stride-pro/backend/internal/auth"
	"github.com/stride-pro/backend/internal/barns"
	biz "github.com/stride-pro/backend/internal/business_settings"
	"github.com/stride-pro/backend/internal/clients"
	"github.com/stride-pro/backend/internal/database"
	"github.com/stride-pro/backend/internal/horses"
	"github.com/stride-pro/backend/internal/invoices"
	"github.com/stride-pro/backend/internal/middleware"
	"github.com/stride-pro/backend/internal/sessions"
	svc "github.com/stride-pro/backend/internal/service_items"
	"github.com/stride-pro/backend/internal/subscriptions"
	"github.com/stride-pro/backend/pkg/response"
)

// Deps holds all handler dependencies needed to configure routing.
type Deps struct {
	DB                     *database.DB
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
}

// New creates and configures the application router with all routes and middleware.
func New(deps Deps) http.Handler {
	r := mux.NewRouter()

	// Global middleware
	corsConfig := middleware.DefaultCORSConfig()
	r.Use(middleware.CORS(corsConfig))
	r.Use(middleware.Logging)

	rateLimiter := middleware.NewRateLimiter(10, 50) // 10 req/s, burst of 50
	r.Use(rateLimiter.Middleware)

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

	// Public auth routes
	r.HandleFunc("/api/auth/register", deps.AuthHandler.Register).Methods("POST")
	r.HandleFunc("/api/auth/login", deps.AuthHandler.Login).Methods("POST")
	r.HandleFunc("/api/auth/refresh", deps.AuthHandler.Refresh).Methods("POST")

	// Protected routes
	protected := r.PathPrefix("/api").Subrouter()
	protected.Use(auth.Middleware(deps.AuthService))

	// Auth - protected
	protected.HandleFunc("/auth/me", deps.AuthHandler.Me).Methods("GET")

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

	// Settings
	protected.HandleFunc("/settings/business", deps.BusinessSettingHandler.Get).Methods("GET")
	protected.HandleFunc("/settings/business", deps.BusinessSettingHandler.Upsert).Methods("PUT")
	protected.HandleFunc("/settings/service-items", deps.ServiceItemHandler.List).Methods("GET")
	protected.HandleFunc("/settings/service-items", deps.ServiceItemHandler.Create).Methods("POST")
	protected.HandleFunc("/settings/service-items/{id}", deps.ServiceItemHandler.Update).Methods("PUT")
	protected.HandleFunc("/settings/service-items/{id}", deps.ServiceItemHandler.Delete).Methods("DELETE")

	return r
}
