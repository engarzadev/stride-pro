import { Injectable, inject } from '@angular/core';
import { Observable } from 'rxjs';
import { Client } from '../../core/models';
import { ApiService } from '../../core/services/api.service';

@Injectable({ providedIn: 'root' })
export class ClientsService {
  private readonly api = inject(ApiService);

  getAll(): Observable<Client[]> {
    return this.api.get<Client[]>('/clients');
  }

  getById(id: string): Observable<Client> {
    return this.api.get<Client>(`/clients/${id}`);
  }

  create(client: Partial<Client>): Observable<Client> {
    return this.api.post<Client>('/clients', this.toPayload(client));
  }

  update(id: string, client: Partial<Client>): Observable<Client> {
    return this.api.put<Client>(`/clients/${id}`, this.toPayload(client));
  }

  delete(id: string): Observable<void> {
    return this.api.delete<void>(`/clients/${id}`);
  }

  private toPayload(client: Partial<Client>): object {
    return {
      first_name: client.firstName,
      last_name: client.lastName,
      email: client.email,
      phone: client.phone,
      address: client.address,
      notes: client.notes,
    };
  }
}
