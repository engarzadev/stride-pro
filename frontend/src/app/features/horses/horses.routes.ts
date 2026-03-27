import { Routes } from '@angular/router';
import { HorseListComponent } from './horse-list/horse-list.component';
import { HorseDetailComponent } from './horse-detail/horse-detail.component';
import { HorseFormComponent } from './horse-form/horse-form.component';

export const horsesRoutes: Routes = [
  { path: '', component: HorseListComponent },
  { path: 'new', component: HorseFormComponent },
  { path: ':id', component: HorseDetailComponent },
  { path: ':id/edit', component: HorseFormComponent },
];
