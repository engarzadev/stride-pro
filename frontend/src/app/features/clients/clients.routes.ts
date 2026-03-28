import { Routes } from '@angular/router';
import { ClientListComponent } from './client-list/client-list.component';
import { ClientDetailComponent } from './client-detail/client-detail.component';
import { ClientFormComponent } from './client-form/client-form.component';
import { ClientOnboardingComponent } from './client-onboarding/client-onboarding.component';

export const clientsRoutes: Routes = [
  { path: '', component: ClientListComponent },
  { path: 'onboard', component: ClientOnboardingComponent },
  { path: 'new', component: ClientFormComponent },
  { path: ':id', component: ClientDetailComponent },
  { path: ':id/edit', component: ClientFormComponent },
];
