# Stride Pro

Full-stack equine management platform for horse owners and equine professionals (trainers, farriers, vets). Manage horses, clients, appointments, care logs, invoicing, and more from a single application.

## Overview

Stride Pro serves two distinct user types with tailored experiences:

- **Horse Owners** — track their own horses' health, care history, and upcoming reminders.
- **Professionals** — manage clients, schedule appointments, record session notes, send invoices, and run their equine business.

## Tech Stack

| Layer | Technology |
|-------|------------|
| Frontend | Angular 19, Angular Material 19, SCSS, TypeScript |
| Backend | Go 1.22, gorilla/mux, golang-jwt |
| Database | PostgreSQL |
| Deployment | Docker, Railway |

## Project Structure

```
stride-pro/
├── frontend/          # Angular 19 SPA
├── backend/           # Go REST API
├── scripts/           # Utility scripts (seeding, etc.)
├── Dockerfile         # Production container (Go API + migrations)
└── railway.json       # Railway deployment config
```

See [frontend/README.md](frontend/README.md) and [backend/README.md](backend/README.md) for detailed architecture documentation.

## Account Types

Users choose an account type at registration. Each type has a tailored dashboard, navigation, and feature set.

### Owner

For individual horse owners managing their own horses.

| Feature | Description |
|---------|-------------|
| Horses | Add and manage horse profiles (breed, age, vet/farrier contacts, notes) |
| Care Log | Record care events — farrier visits, vaccinations, dental, diet changes, and more |
| Reminders | Manual reminders + auto-generated reminders based on care log history |
| Dashboard | At-a-glance view of horses and upcoming reminders |

### Professional

For trainers, farriers, vets, and equine service providers.

| Feature | Description |
|---------|-------------|
| Clients | Manage horse owner contacts and their associated horses |
| Horses | Full horse records with client and barn assignment |
| Barns | Manage barn locations where horses are kept |
| Appointments | Schedule visits with calendar view, status tracking, and travel time |
| Sessions | Record detailed session notes, findings, and recommendations per appointment |
| Invoices | Create and manage invoices with line items and service catalog |
| Dashboard | Today's appointments, recent activity, and business summary |

Both account types have access to **Settings** (account profile, business settings for professionals) and **Billing** (subscription management).

## Subscription Plans

All users start on the **Free** plan. Limits are enforced at the API layer — requests that exceed a plan's limits return `403 Forbidden`. The frontend displays upgrade banners for gated features.

### Owner Plans

| Feature | Free | Base ($29.99/mo) |
|---------|------|------------------|
| Horses | Up to 20 | Unlimited |
| Care log | Reminders only | Full (log + reminders) |

### Professional Plans

| Feature | Free | Base ($29.99/mo) | Trainer ($49.99/mo) | Enterprise ($99.99/mo) |
|---------|------|------------------|---------------------|------------------------|
| Clients | Up to 10 | Unlimited | Unlimited | Unlimited |
| Horses | Up to 20 | Unlimited | Unlimited | Unlimited |
| Appointments | Basic | Full | Full | Full |
| Invoices | Basic | Full (with templates) | Full | Full |
| Session notes | — | Yes | Yes | Yes |
| Barn management | — | Yes | Yes | Yes |
| Care logs | — | Yes | Yes | Yes |
| Email notifications | — | Yes | Yes | Yes |
| SMS notifications | — | — | Yes | Yes |
| Multi-horse sessions | — | — | Yes | Yes |
| Advanced reporting | — | — | Yes | Yes |
| Client portal | — | — | Yes | Yes |
| API access | — | — | — | Yes |
| Custom branding | — | — | — | Yes |
| Priority support | — | — | — | Yes |

## Getting Started

### Prerequisites

- Node.js 18+
- Go 1.22+
- PostgreSQL 15+

### Setup

```bash
# Backend
cd backend
cp .env.example .env          # Set DATABASE_URL and JWT_SECRET
go mod download
go run ./cmd/api              # Starts on :8080

# Frontend (in a separate terminal)
cd frontend
npm install
npm start                     # Starts on :4200, proxies /api to :8080
```

### Environment Variables

See [`backend/.env.example`](backend/.env.example) for the full list. Required for local development:

| Variable | Description |
|----------|-------------|
| `DATABASE_URL` | PostgreSQL connection string |
| `JWT_SECRET` | JWT signing secret (min 32 characters) |

Optional integrations:

| Variable | Description |
|----------|-------------|
| `SENDGRID_API_KEY` | SendGrid API key for email notifications |
| `TWILIO_ACCOUNT_SID` / `TWILIO_AUTH_TOKEN` / `TWILIO_PHONE_NUMBER` | Twilio credentials for SMS notifications |
| `TLS_CERT_FILE` / `TLS_KEY_FILE` | Direct TLS termination |
| `TLS_PROXY_MODE` | Set `true` when behind a TLS-terminating proxy |

### Deployment

The project ships as a single Docker container (Go API serving the Angular build). Railway deployment is configured via `railway.json` and `Dockerfile`.

```bash
# Build the Docker image
docker build -t stride-pro .
```
