// Package main is the entry point for the Stride Pro API server.
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/stride-pro/backend/internal/appointments"
	"github.com/stride-pro/backend/internal/auth"
	"github.com/stride-pro/backend/internal/barns"
	biz "github.com/stride-pro/backend/internal/business_settings"
	"github.com/stride-pro/backend/internal/clients"
	"github.com/stride-pro/backend/internal/config"
	"github.com/stride-pro/backend/internal/database"
	"github.com/stride-pro/backend/internal/horses"
	"github.com/stride-pro/backend/internal/invoices"
	"github.com/stride-pro/backend/internal/notifications"
	"github.com/stride-pro/backend/internal/router"
	"github.com/stride-pro/backend/internal/sessions"
	svc "github.com/stride-pro/backend/internal/service_items"
	"github.com/stride-pro/backend/internal/subscriptions"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Connect to database
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()
	log.Println("connected to database")

	// Run migrations
	if err := db.RunMigrations("migrations"); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}
	log.Println("migrations applied")

	// Initialize services
	authService := auth.NewService(db, cfg.JWTSecret)
	authHandler := auth.NewHandler(authService)

	subsService := subscriptions.NewService(db)
	subsHandler := subscriptions.NewHandler(subsService)

	clientRepo := clients.NewRepository(db)
	clientService := clients.NewService(clientRepo, subsService)
	clientHandler := clients.NewHandler(clientService)

	horseRepo := horses.NewRepository(db)
	horseService := horses.NewService(horseRepo, subsService)
	horseHandler := horses.NewHandler(horseService)

	barnRepo := barns.NewRepository(db)
	barnService := barns.NewService(barnRepo, subsService)
	barnHandler := barns.NewHandler(barnService)

	var emailSender notifications.EmailSender
	if cfg.SendGridAPIKey != "" {
		emailSender = notifications.NewSendGridEmailSender(cfg.SendGridAPIKey, cfg.SendGridFromEmail, cfg.SendGridFromName)
		log.Println("email: using SendGrid")
	} else {
		emailSender = notifications.NewStubEmailSender()
		log.Println("email: using stub (set SENDGRID_API_KEY to enable real sending)")
	}
	notifService := notifications.NewService(db, emailSender, notifications.NewStubSMSSender())

	apptRepo := appointments.NewRepository(db)
	apptService := appointments.NewService(apptRepo, notifService, authService)
	apptHandler := appointments.NewHandler(apptService)

	sessionRepo := sessions.NewRepository(db)
	sessionService := sessions.NewService(sessionRepo, subsService)
	sessionHandler := sessions.NewHandler(sessionService)

	bizRepo := biz.NewRepository(db)
	bizService := biz.NewService(bizRepo)
	bizHandler := biz.NewHandler(bizService)

	svcRepo := svc.NewRepository(db)
	svcService := svc.NewService(svcRepo)
	svcHandler := svc.NewHandler(svcService)

	invoiceRepo := invoices.NewRepository(db)
	invoiceService := invoices.NewService(invoiceRepo, bizService, emailSender)
	invoiceHandler := invoices.NewHandler(invoiceService)

	// Set up router
	handler := router.New(router.Deps{
		DB:                     db,
		AuthService:            authService,
		AuthHandler:            authHandler,
		ClientHandler:          clientHandler,
		HorseHandler:           horseHandler,
		BarnHandler:            barnHandler,
		ApptHandler:            apptHandler,
		SessionHandler:         sessionHandler,
		InvoiceHandler:         invoiceHandler,
		SubscriptionHandler:    subsHandler,
		BusinessSettingHandler: bizHandler,
		ServiceItemHandler:     svcHandler,
	})

	// Configure HTTP server
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("server starting on port %s (env: %s)", cfg.Port, cfg.Environment)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

	log.Println("server stopped")
}
