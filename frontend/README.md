# Stride Pro Frontend

Angular 19 single-page application for the Stride Pro equine management platform.

## Architecture

The frontend follows a **Core / Features / Shared** architecture with lazy-loaded feature modules and standalone components.

```
src/app/
├── core/        # Singleton services, interceptors, guards, models, utils
├── features/    # Lazy-loaded feature domains (one folder per domain)
└── shared/      # Reusable components, directives, and pipes
```

### Core

The `core/` layer is loaded once at startup. Nothing outside of `core/` should depend on it except `features/`.

| Path | Purpose |
|------|---------|
| `services/api.service.ts` | Generic HTTP wrapper (`get`, `post`, `put`, `patch`, `delete`) — all feature services use this instead of `HttpClient` directly. Handles camelCase ↔ snake_case conversion and unwraps the `{ data }` response envelope. |
| `services/auth.service.ts` | Authentication state — login, register, logout. Persists user to `localStorage`. |
| `services/subscription.service.ts` | Loads the current user's subscription plan and exposes `hasFeature()` for feature gating. |
| `services/theme.service.ts` | Dark/light mode toggle — adds/removes `dark` class on `<html>`. |
| `interceptors/auth.interceptor.ts` | Sets `withCredentials: true` on all requests so HttpOnly cookies are sent. |
| `interceptors/csrf.interceptor.ts` | Reads the `XSRF-TOKEN` cookie and attaches it as the `X-XSRF-TOKEN` header on mutating requests. |
| `interceptors/error.interceptor.ts` | Centralized HTTP error handling — shows toast messages, handles 401 session expiry redirect, suppresses subscription-gated 403s. |
| `guards/auth.guard.ts` | Redirects unauthenticated users to `/auth/login`. |
| `models/index.ts` | Shared TypeScript interfaces (`User`, `Client`, `Horse`, `Barn`, `Appointment`, `Session`, `Invoice`, `Reminder`, `CareLog`, `SubscriptionPlan`, etc.) |
| `utils/camel-case.ts` | `keysToCamel` / `keysToSnake` recursive object key converters used by `ApiService`. |

### Features

Each domain under `features/` is self-contained and lazy-loaded. The typical structure inside a feature folder is:

```
features/<domain>/
├── <domain>-list/       # List view component
├── <domain>-detail/     # Read-only detail view component
├── <domain>-form/       # Create / edit form component
├── <domain>.service.ts  # HTTP calls scoped to this domain
└── <domain>.routes.ts   # Route definitions (lazy-loaded via loadChildren)
```

| Feature | Domain | Account Type |
|---------|--------|--------------|
| `auth/` | Login and registration | Both |
| `dashboard/` | Summary/overview page | Both |
| `horses/` | Horse records, care log, reminders | Both |
| `owner-care-log/` | Owner-facing care log page | Owner |
| `owner-reminders/` | Owner-facing reminders page | Owner |
| `clients/` | Client management | Professional |
| `barns/` | Barn management | Professional |
| `appointments/` | Appointment scheduling | Professional |
| `sessions/` | Training session tracking | Professional |
| `invoices/` | Invoice management | Professional |
| `billing/` | Subscription plan management | Both |
| `settings/` | Account and business settings | Both |

### Shared

`shared/` contains UI building blocks with no business logic.

**Components:**

| Component | Purpose |
|-----------|---------|
| `HeaderComponent` | Top navigation bar with user menu |
| `SidebarComponent` | Side navigation — items filtered by account type (owner vs professional) |
| `BottomNavComponent` | Mobile bottom navigation bar with quick-add menu |
| `PageHeaderComponent` | Consistent page title with optional action button |
| `FormPageComponent` | Page layout for create/edit forms with save/cancel actions |
| `DetailPageComponent` | Page layout for read-only detail views |
| `DataTableComponent` | Reusable sortable data table with configurable columns and row actions |
| `LoadingSpinnerComponent` | Loading indicator |
| `ConfirmDialogComponent` | Confirmation modal |
| `ToastService` | Notification toasts (success, error) |
| `BreedAutocompleteComponent` | Searchable horse breed selector |
| `HorseMultiselectAutocompleteComponent` | Multi-horse selector for sessions |
| `QuickCreateBarnComponent` / `QuickCreateClientComponent` | Inline create dialogs from within forms |
| `UpgradeFieldPromptComponent` | Inline upgrade prompt for gated form fields |

**Pipes:**

| Pipe | Purpose |
|------|---------|
| `DateFormatPipe` | Consistent date formatting |
| `CurrencyFormatPipe` | Currency display formatting |

**Directives:**

| Directive | Purpose |
|-----------|---------|
| `ClickOutsideDirective` | Detects clicks outside a host element |

## Key Patterns

### Standalone Components
All components use `standalone: true`. Imports are declared per-component rather than in NgModules.

### Dependency Injection
Uses the modern `inject()` function instead of constructor injection:
```typescript
private readonly clientsService = inject(ClientsService);
private readonly router = inject(Router);
```

### Reactive State with Signals
Component-local state uses Angular signals; async data from services uses RxJS observables:
```typescript
clients = signal<Client[]>([]);
isLoading = signal(false);
```

### Feature Service Pattern
Every feature service follows the same CRUD contract through `ApiService`:
```typescript
getAll()         → Observable<T[]>
getById(id)      → Observable<T>
create(body)     → Observable<T>
update(id, body) → Observable<T>
delete(id)       → Observable<void>
```

### API Layer
All HTTP calls go through `ApiService`. It:
1. Converts request bodies from camelCase → snake_case before sending
2. Unwraps the `{ data: ... }` response envelope
3. Converts response keys from snake_case → camelCase

Never use `HttpClient` directly — always go through `ApiService`.

### Theming
Dark mode is implemented via CSS custom properties. `ThemeService` toggles a `dark` class on `<html>`. Component styles use `var(--color-*)` tokens defined in `styles.scss`. For dark mode overrides in component SCSS, use `:host-context(html.dark)`.

### Routing
Routes are lazy-loaded via `loadChildren` / `loadComponent` in `app.routes.ts`. All routes except `/auth` are protected by `authGuard`. The root path redirects to `/dashboard`.

## Project Structure

```
frontend/
├── src/
│   ├── app/
│   │   ├── core/
│   │   │   ├── guards/
│   │   │   ├── interceptors/
│   │   │   ├── models/
│   │   │   ├── services/
│   │   │   └── utils/
│   │   ├── features/
│   │   ├── shared/
│   │   │   ├── components/
│   │   │   ├── directives/
│   │   │   └── pipes/
│   │   ├── app.component.ts
│   │   ├── app.config.ts
│   │   └── app.routes.ts
│   ├── environments/
│   │   ├── environment.ts          # Dev: /api (proxied to localhost:8080)
│   │   └── environment.prod.ts     # Prod: https://stride-pro-production.up.railway.app/api
│   ├── styles.scss                 # Global styles + CSS custom properties
│   └── main.ts
├── proxy.conf.json                 # Dev proxy: /api → localhost:8080
├── angular.json
├── tsconfig.json
└── package.json
```

## Running

```bash
npm install          # Install dependencies
npm start            # Dev server at http://localhost:4200 (proxies /api to backend)
npm run build        # Production build → dist/stride-pro/
ng test              # Run unit tests
```

The dev server proxies `/api` requests to the Go backend at `http://localhost:8080` via `proxy.conf.json`.
