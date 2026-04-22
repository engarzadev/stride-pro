# Stride Pro Backend

Go REST API powering the Stride Pro equine management platform.

## Architecture

The backend follows a **layered architecture** pattern. Each feature domain is organized into three layers:

### Handler → Service → Repository

```
HTTP Request → handler → service → repository → database
HTTP Response ← handler ← service ← repository ← database
```

| Layer | File | Responsibility |
|-------|------|----------------|
| **Handler** | `handler.go` | Parses HTTP requests (URL params, JSON body), calls the service layer, and writes HTTP responses (JSON, status codes). This is the only layer that knows about `http.Request` and `http.ResponseWriter`. |
| **Service** | `service.go` | Contains business logic, validation, and authorization rules. Orchestrates calls to one or more repositories. Knows nothing about HTTP. |
| **Repository** | `repository.go` | Handles all database access (SQL queries, inserts, updates, deletes). Returns Go structs. Knows nothing about HTTP or business rules. |

### Why this pattern?

- **Separation of concerns** — each layer has a single responsibility.
- **Testability** — services can be tested without HTTP, repositories can be tested against a test database.
- **Flexibility** — swapping the database only affects the repository layer. Changing a business rule only affects the service layer.

## Project Structure

```
backend/
├── cmd/
│   └── api/
│       └── main.go              # Application entrypoint
├── internal/
│   ├── appointments/            # Scheduling and calendar management
│   │   ├── handler.go
│   │   ├── service.go
│   │   └── repository.go
│   ├── auth/                    # Authentication (login, register, JWT, middleware)
│   ├── barns/                   # Barn / location management
│   ├── business_settings/       # Business profile for invoicing
│   ├── care_logs/               # Horse care event tracking (farrier, vet, diet, etc.)
│   ├── clients/                 # Client (horse-owner) management
│   ├── horses/                  # Horse records and details
│   ├── invoices/                # Invoicing and line items
│   ├── notifications/           # Email (SendGrid) and SMS (Twilio) notifications
│   ├── reminders/               # Upcoming care reminders (manual and auto-generated)
│   ├── service_items/           # Reusable service price catalog
│   ├── sessions/                # Training / care session notes
│   ├── subscriptions/           # Subscription plans and feature gating
│   ├── config/                  # Environment variable configuration
│   ├── database/                # Database connection and health check
│   ├── middleware/              # CORS, CSRF, logging, rate limiting, recovery, security headers
│   ├── models/                  # Shared data models
│   └── router/                  # Route definitions
├── migrations/                  # SQL migration files
├── pkg/
│   ├── response/                # Standardized JSON response helpers
│   └── validator/               # Input validation utilities
├── go.mod
└── go.sum
```

## Subscription Feature Gating

All new users start on the **Free** plan. Plan limits are enforced at the service layer — requests that exceed a limit return `403 Forbidden`. See the [root README](../README.md) for account types and subscription plan details.

The account type is stored as the `role` field on the user record (`owner` or `professional`). The subscription tier is stored in `users.subscription_tier`.

### Feature Flags

Plans are defined in `internal/subscriptions/model.go`. Each plan has a list of feature flags checked at the service layer via `SubscriptionService.RequireFeature()` or `HasFeature()`. Key flags:

| Flag | Description |
|------|-------------|
| `clients_max_10` / `clients_unlimited` | Client resource limits |
| `horses_max_20` / `horses_unlimited` | Horse resource limits |
| `care_logs` | Full care log access |
| `care_log_reminders` | Auto-generated care reminders |
| `session_notes` | Session notes and findings |
| `barn_management` | Barn / location management |
| `email_notifications` | Email via SendGrid |
| `sms_notifications` | SMS via Twilio |
| `multi_horse_sessions` | Multi-horse session tracking |
| `advanced_reporting` | Advanced reporting and analytics |
| `client_portal` | Client self-service portal |
| `api_access` | REST API access |
| `custom_branding` | Custom branding |
| `priority_support` | Priority support |

The `GET /api/subscription` endpoint returns the current plan and live usage counts for enforced resources.

## Running

```bash
cp .env.example .env   # Configure database URL and JWT secret
go mod download         # Download dependencies
go run ./cmd/api        # Start the server on :8080
```

## Configuration

See `.env.example` for all available environment variables. See the [root README](../README.md#environment-variables) for a summary of required and optional variables.
