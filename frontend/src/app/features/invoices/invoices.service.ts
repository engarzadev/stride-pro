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

  getById(id: string): Observable<Invoice> {
    return this.api.get<Invoice>(`/invoices/${id}`);
  }

  create(invoice: InvoicePayload): Observable<Invoice> {
    return this.api.post<Invoice>('/invoices', invoice);
  }

  update(id: string, invoice: InvoicePayload): Observable<Invoice> {
    return this.api.put<Invoice>(`/invoices/${id}`, invoice);
  }

  delete(id: string): Observable<void> {
    return this.api.delete<void>(`/invoices/${id}`);
  }

  sendInvoice(id: string): Observable<{ message: string }> {
    return this.api.post<{ message: string }>(`/invoices/${id}/send`, {});
  }

  updateStatus(id: string, status: string): Observable<{ message: string }> {
    return this.api.patch<{ message: string }>(`/invoices/${id}/status`, { status });
  }
}
