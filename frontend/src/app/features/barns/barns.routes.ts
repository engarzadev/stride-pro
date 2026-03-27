import { Routes } from '@angular/router';
import { BarnListComponent } from './barn-list/barn-list.component';
import { BarnDetailComponent } from './barn-detail/barn-detail.component';
import { BarnFormComponent } from './barn-form/barn-form.component';

export const barnsRoutes: Routes = [
  { path: '', component: BarnListComponent },
  { path: 'new', component: BarnFormComponent },
  { path: ':id', component: BarnDetailComponent },
  { path: ':id/edit', component: BarnFormComponent },
];
