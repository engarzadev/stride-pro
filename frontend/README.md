# Stride Pro Frontend

Angular 19 single-page application for the Stride Pro equine management platform.

## Architecture

The frontend follows a **Core / Features / Shared** architecture with lazy-loaded feature modules and standalone components.

```
src/app/
├── core/        # Singleton services, interceptors, guards, models
├── features/    # Lazy-loaded feature domains (one folder per domain)
└── shared/      # Reusable components, directives, and pipes
```

### Core

The `core/` layer is loaded once at startup. Nothing outside of `core/` should depend on it except `features/`.

| Path | Purpose |
|------|---------|
| `core/services/api.service.ts` | Generic HTTP wrapper (`get`, `post`, `put`, `delete`) — all feature services use this instead of `HttpClient` directly |
| `core/services/auth.service.ts` | Authentication state — exposes `isAuthenticated$` and `currentUser$` as observables backed by `BehaviorSubject`; persists token and user to `localStorage` |
| `core/interceptors/auth.interceptor.ts` | Attaches `Authorization: Bearer <token>` to every outgoing request |
| `core/interceptors/error.interceptor.ts` | Centralized HTTP error handling |
| `core/guards/auth.guard.ts` | Redirects unauthenticated users away from protected routes |
| `core/models/index.ts` | Shared TypeScript interfaces (`User`, `Client`, `Horse`, `Barn`, `Appointment`, `Session`, `Invoice`, `ApiResponse<T>`, `PaginatedResponse<T>`) |

### Features

Each domain under `features/` is self-contained and la`zy-loaded. The typical structure inside a feature folder is:

```
features/<domain>/
├── <domain>-list/       # List view component
├── <domain>-detail/     # Read-only detail view component
├── <domain>-form/       # Create / edit form component
├── <domain>.service.ts  # HTTP calls scoped to this domain
└── <domain>.routes.ts   # Route definitions (lazy-loaded via loadComponent)
```

| Feature | Domain |
|---------|--------|
| `auth/` | Login and registration pages |
| `dashboard/` | Summary/overview page |
| `clients/` | Client management |
| `horses/` | Horse management |
| `barns/` | Barn management |
| `appointments/` | Appointment scheduling |
| `sessions/` | Training session tracking |
| `invoices/` | Invoice management |

### Shared

`shared/` contains UI building blocks with no business logic.

| Component / Service | Purpose |
|---------------------|---------|
| `HeaderComponent` | Top navigation bar |
| `SidebarComponent` | Side navigation |
| `PageHeaderComponent` | Consistent page title area |
| `DataTableComponent` | Reusable sortable data table with configurable columns and row actions |
| `LoadingSpinnerComponent` | Full-screen loading indicator |
| `ConfirmDialogService` | Imperative confirmation modal |
| `ToastService` | Notification toasts |

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
getAll()         → Observable<PaginatedResponse<T>>
getById(id)      → Observable<T>
create(body)     → Observable<T>
update(id, body) → Observable<T>
delete(id)       → Observable<void>
```

### Routing
Routes are lazy-loaded via `loadComponent`. Protected routes use `canActivate: [authGuard]`. The root path redirects to `/dashboard`.

## Project Structure

```
frontend/
├── src/
│   ├── app/
│   │   ├── core/
│   │   ├── features/
│   │   ├── shared/
│   │   ├── app.component.ts
│   │   └── app.routes.ts
│   ├── environments/
│   │   ├── environment.ts          # Dev: http://localhost:8080/api
│   │   └── environment.prod.ts     # Prod API URL
│   ├── styles.scss                 # Global styles
│   └── main.ts
├── angular.json
├── tsconfig.json
└── package.json
```

## Running

```bash
npm install          # Install dependencies
npm start            # Dev server at http://localhost:4200
npm run build        # Production build → dist/stride-pro/
npm test             # Run unit tests
```

The dev server proxies API calls to the Go backend at `http://localhost:8080/api` (configured in `src/environments/environment.ts`).
