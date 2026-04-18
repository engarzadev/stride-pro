import { Routes } from '@angular/router';
import { SessionListComponent } from './session-list/session-list.component';
import { SessionDetailComponent } from './session-detail/session-detail.component';
import { SessionFormComponent } from './session-form/session-form.component';

export const sessionsRoutes: Routes = [
  { path: '', component: SessionListComponent },
  { path: 'new', component: SessionFormComponent },
  { path: ':id', component: SessionDetailComponent },
  { path: ':id/edit', component: SessionFormComponent },
];
