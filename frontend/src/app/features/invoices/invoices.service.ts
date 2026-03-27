import { Injectable, inject } from '@angular/core';
import { Observable } from 'rxjs';
import { ApiService } from '../../core/services/api.service';
import { Invoice, InvoiceItem } from '../../core/models';

type InvoicePayload = Omit<Partial<Invoice>, 'items'> & {
  items?: Omit<InvoiceItem, 'id' | 'invoiceId'>[];
};

@Injectable({ providedIn: 'root' })
export class InvoicesService {
  private readonly api = inject(ApiService);

  getAll(): Observable<Invoice[]> {
    return this.api.get<Invoice[]>('/invoices');
  }

  getById(id: number): Observable<Invoice> {
    return this.api.get<Invoice>(`/invoices/${id}`);
  }

  create(invoice: InvoicePayload): Observable<Invoice> {
    return this.api.post<Invoice>('/invoices', invoice);
  }

  update(id: number, invoice: InvoicePayload): Observable<Invoice> {
    return this.api.put<Invoice>(`/invoices/${id}`, invoice);
  }

  delete(id: number): Observable<void> {
    return this.api.delete<void>(`/invoices/${id}`);
  }
}
