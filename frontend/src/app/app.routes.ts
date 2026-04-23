import { Routes } from '@angular/router';
import { authGuard } from './core/guards/auth.guard';
import { subscriptionResolver } from './core/resolvers/subscription.resolver';

const protectedRoute = (path: string, opts: object) => ({
  path,
  canActivate: [authGuard],
  resolve: { subscription: subscriptionResolver },
  ...opts,
});

export const routes: Routes = [
  { path: '', redirectTo: 'dashboard', pathMatch: 'full' },
  {
    path: 'auth',
    loadChildren: () => import('./features/auth/auth.routes').then((m) => m.authRoutes),
  },
  protectedRoute('dashboard', {
    loadChildren: () => import('./features/dashboard/dashboard.routes').then((m) => m.dashboardRoutes),
  }),
  protectedRoute('clients', {
    loadChildren: () => import('./features/clients/clients.routes').then((m) => m.clientsRoutes),
  }),
  protectedRoute('horses', {
    loadChildren: () => import('./features/horses/horses.routes').then((m) => m.horsesRoutes),
  }),
  protectedRoute('care-log', {
    loadComponent: () => import('./features/owner-care-log/owner-care-log.component').then((m) => m.OwnerCareLogComponent),
  }),
  protectedRoute('reminders', {
    loadComponent: () => import('./features/owner-reminders/owner-reminders.component').then((m) => m.OwnerRemindersComponent),
  }),
  protectedRoute('barns', {
    loadChildren: () => import('./features/barns/barns.routes').then((m) => m.barnsRoutes),
  }),
  protectedRoute('appointments', {
    loadChildren: () => import('./features/appointments/appointments.routes').then((m) => m.appointmentsRoutes),
  }),
  protectedRoute('sessions', {
    loadChildren: () => import('./features/sessions/sessions.routes').then((m) => m.sessionsRoutes),
  }),
  protectedRoute('invoices', {
    loadChildren: () => import('./features/invoices/invoices.routes').then((m) => m.invoicesRoutes),
  }),
  protectedRoute('billing', {
    loadChildren: () => import('./features/billing/billing.routes').then((m) => m.billingRoutes),
  }),
  protectedRoute('settings', {
    loadChildren: () => import('./features/settings/settings.routes').then((m) => m.settingsRoutes),
  }),
  { path: '**', redirectTo: 'dashboard' },
];
