import { Routes } from '@angular/router';
import { AppointmentListComponent } from './appointment-list/appointment-list.component';
import { AppointmentDetailComponent } from './appointment-detail/appointment-detail.component';
import { AppointmentFormComponent } from './appointment-form/appointment-form.component';

export const appointmentsRoutes: Routes = [
  { path: '', component: AppointmentListComponent },
  { path: 'new', component: AppointmentFormComponent },
  { path: ':id', component: AppointmentDetailComponent },
  { path: ':id/edit', component: AppointmentFormComponent },
];
