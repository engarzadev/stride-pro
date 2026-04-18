# Stride Pro Backend

Go REST API powering the Stride Pro equine management platform.

## Architecture

The backend follows a **layered architecture** pattern. Each feature domain (appointments, barns, clients, horses, sessions, invoices) is organized into three layers:

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
│   ├── appointments/            # Appointments feature
│   │   ├── handler.go
│   │   ├── service.go
│   │   └── repository.go
│   ├── auth/                    # Authentication (handler, middleware, service)
│   ├── barns/                   # Barns feature
│   ├── clients/                 # Clients feature
│   ├── horses/                  # Horses feature
│   ├── sessions/                # Training sessions feature
│   ├── invoices/                # Invoices feature
│   ├── notifications/           # Email and SMS notifications
│   ├── subscriptions/           # Subscription management
│   ├── config/                  # Environment variable configuration
│   ├── database/                # Database connection and migrations
│   ├── middleware/               # CORS, logging, rate limiting
│   ├── models/                  # Shared data models
│   └── router/                  # Route definitions
├── migrations/                  # SQL migration files
├── pkg/
│   ├── response/                # Standardized JSON response helpers
│   └── validator/               # Input validation utilities
├── go.mod
└── go.sum
```

## Running

```bash
cp .env.example .env   # Configure database URL and JWT secret
go mod download         # Download dependencies
go run ./cmd/api        # Start the server on :8080
```

## Subscription Plans

All new users start on the **Free** plan. Plan limits are enforced at the API layer — requests that exceed a limit return `403 Forbidden`.

| Feature | Free | Base ($29.99) | Trainer Add-on ($49.99) | Enterprise ($99.99) |
|---------|------|---------------|-------------------------|---------------------|
| Clients | Max 10 | Unlimited | Unlimited | Unlimited |
| Horses | Max 20 | Unlimited | Unlimited | Unlimited |
| Appointments | Basic | Full (with reminders) | Full | Full |
| Invoices | Basic | Full (with templates) | Full | Full |
| Session notes | — | Yes | Yes | Yes |
| Barn management | — | Yes | Yes | Yes |
| Email notifications | — | Yes | Yes | Yes |
| SMS notifications | — | — | Yes | Yes |
| Multi-horse sessions | — | — | Yes | Yes |
| Advanced reporting | — | — | Yes | Yes |
| Client portal | — | — | Yes | Yes |
| API access | — | — | — | Yes |
| Custom branding | — | — | — | Yes |
| Priority support | — | — | — | Yes |

A user's plan is stored in `users.subscription_tier`. The `GET /api/subscription` endpoint returns the current plan and live usage counts for enforced resources.

## Configuration

See `.env.example` for all available environment variables. Required for local development:

- `DATABASE_URL` — PostgreSQL connection string
- `JWT_SECRET` — Secret key for signing JWT tokens
