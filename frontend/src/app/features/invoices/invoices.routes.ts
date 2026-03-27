import { Routes } from '@angular/router';
import { InvoiceListComponent } from './invoice-list/invoice-list.component';
import { InvoiceDetailComponent } from './invoice-detail/invoice-detail.component';
import { InvoiceFormComponent } from './invoice-form/invoice-form.component';

export const invoicesRoutes: Routes = [
  { path: '', component: InvoiceListComponent },
  { path: 'new', component: InvoiceFormComponent },
  { path: ':id', component: InvoiceDetailComponent },
  { path: ':id/edit', component: InvoiceFormComponent },
];
