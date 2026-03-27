import { Routes } from '@angular/router';
import { ClientListComponent } from './client-list/client-list.component';
import { ClientDetailComponent } from './client-detail/client-detail.component';
import { ClientFormComponent } from './client-form/client-form.component';

export const clientsRoutes: Routes = [
  { path: '', component: ClientListComponent },
  { path: 'new', component: ClientFormComponent },
  { path: ':id', component: ClientDetailComponent },
  { path: ':id/edit', component: ClientFormComponent },
];
