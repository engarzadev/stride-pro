import { Routes } from '@angular/router';
import { authGuard } from './core/guards/auth.guard';

export const routes: Routes = [
  { path: '', redirectTo: 'dashboard', pathMatch: 'full' },
  {
    path: 'auth',
    loadChildren: () => import('./features/auth/auth.routes').then((m) => m.authRoutes),
  },
  {
    path: 'dashboard',
    canActivate: [authGuard],
    loadChildren: () => import('./features/dashboard/dashboard.routes').then((m) => m.dashboardRoutes),
  },
  {
    path: 'clients',
    canActivate: [authGuard],
    loadChildren: () => import('./features/clients/clients.routes').then((m) => m.clientsRoutes),
  },
  {
    path: 'horses',
    canActivate: [authGuard],
    loadChildren: () => import('./features/horses/horses.routes').then((m) => m.horsesRoutes),
  },
  {
    path: 'care-log',
    canActivate: [authGuard],
    loadComponent: () => import('./features/owner-care-log/owner-care-log.component').then((m) => m.OwnerCareLogComponent),
  },
  {
    path: 'reminders',
    canActivate: [authGuard],
    loadComponent: () => import('./features/owner-reminders/owner-reminders.component').then((m) => m.OwnerRemindersComponent),
  },
  {
    path: 'barns',
    canActivate: [authGuard],
    loadChildren: () => import('./features/barns/barns.routes').then((m) => m.barnsRoutes),
  },
  {
    path: 'appointments',
    canActivate: [authGuard],
    loadChildren: () => import('./features/appointments/appointments.routes').then((m) => m.appointmentsRoutes),
  },
  {
    path: 'sessions',
    canActivate: [authGuard],
    loadChildren: () => import('./features/sessions/sessions.routes').then((m) => m.sessionsRoutes),
  },
  {
    path: 'invoices',
    canActivate: [authGuard],
    loadChildren: () => import('./features/invoices/invoices.routes').then((m) => m.invoicesRoutes),
  },
  {
    path: 'billing',
    canActivate: [authGuard],
    loadChildren: () => import('./features/billing/billing.routes').then((m) => m.billingRoutes),
  },
  {
    path: 'settings',
    canActivate: [authGuard],
    loadChildren: () => import('./features/settings/settings.routes').then((m) => m.settingsRoutes),
  },
  { path: '**', redirectTo: 'dashboard' },
];
