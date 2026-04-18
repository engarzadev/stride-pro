# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Stride Pro is a full-stack equine management platform for trainers. It has two separate sub-projects:

- **`frontend/`** — Angular 19 SPA
- **`backend/`** — Go REST API with PostgreSQL

---

## Commands

### Frontend (`cd frontend`)
```bash
npm start          # Dev server at http://localhost:4200
npm run build      # Production build
ng test            # Unit tests (Karma — no spec files exist yet)
ng build --watch   # Watch mode build
```

### Backend (`cd backend`)
```bash
go run ./cmd/api         # Start API server (default :8080)
go build ./...           # Compile all packages
go test ./...            # Run all tests
go test ./internal/...   # Run internal package tests
```

Backend requires a `.env` file at `backend/`. Copy from `.env.example` and set at minimum:
- `DATABASE_URL` — PostgreSQL connection string
- `JWT_SECRET` — JWT signing secret

---

## Architecture

### Frontend

**Directory layout:**
```
src/app/
  core/          # Singletons: services, guards, interceptors, models, utils
  features/      # Lazy-loaded domain modules (auth, dashboard, clients, horses, barns,
                 #   appointments, sessions, invoices, billing, settings)
  shared/        # Reusable components (header, sidebar, data-table, page-header, etc.)
```

**Key conventions:**
- All components are standalone; use `inject()` (not constructor injection) throughout
- Feature routes are lazy-loaded via `loadChildren` in `app.routes.ts`; all routes except `/auth` are protected by `authGuard`
- Use Angular signals (`signal()`, `computed()`) for local component state

**API layer (`core/services/api.service.ts`):**
All HTTP calls go through `ApiService`. It:
1. Automatically converts request bodies from camelCase → snake_case before sending
2. Unwraps the `{ data: ... }` response envelope from the backend
3. Converts response keys from snake_case → camelCase

Never use `HttpClient` directly in components or feature services — always use `ApiService`.

**Theming (`styles.scss` + `core/services/theme.service.ts`):**
Dark mode is implemented via CSS custom properties. `ThemeService` toggles a `dark` class on `<html>`. All component styles should use `var(--color-*)` tokens instead of hardcoded hex values. The available tokens are defined as `:root` and `html.dark` blocks near the top of `styles.scss`.

Styling uses SCSS + Angular Material v19 (MDC-based). The primary palette is forest green (`#3b6255`). For dark mode overrides in component SCSS files, use `:host-context(html.dark)`.

### Backend

**Directory layout:**
```
backend/
  cmd/api/main.go          # Entry point — wires dependencies
  internal/
    <domain>/              # e.g. clients/, horses/, appointments/
      handler.go           # HTTP handlers
      service.go           # Business logic
      repository.go        # SQL queries
    middleware/            # CORS, auth, logging, rate limiting
    router/router.go       # All route definitions
    database/              # DB connection + HealthCheck
    models/                # Shared model types
  migrations/              # SQL migration files
  pkg/response/            # Standardized JSON response helpers
  pkg/validator/           # Input validation utilities
```

**Key conventions:**
- All JSON responses are wrapped: `{ "data": ... }` for success, `{ "error": "..." }` for errors
- All DB columns are snake_case; Go structs use camelCase with `json` tags
- Each domain package follows the Handler → Service → Repository pattern; handlers never touch the DB directly
- `auth.Middleware` validates JWT and sets `userID` on the request context; handlers retrieve it via the context

**Subscription tiers** (enforced at the API service layer, not in handlers):
- `free`: max 10 clients, 20 horses, basic features
- `base` ($29.99): unlimited clients/horses
- `trainer` ($49.99): multi-horse sessions, advanced reporting
- `enterprise` ($99.99): SMS, API access, custom branding
