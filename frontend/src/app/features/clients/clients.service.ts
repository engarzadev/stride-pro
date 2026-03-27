import { Injectable, inject } from '@angular/core';
import { Observable } from 'rxjs';
import { ApiService } from '../../core/services/api.service';
import { Client } from '../../core/models';

@Injectable({ providedIn: 'root' })
export class ClientsService {
  private readonly api = inject(ApiService);

  getAll(): Observable<Client[]> {
    return this.api.get<Client[]>('/clients');
  }

  getById(id: number): Observable<Client> {
    return this.api.get<Client>(`/clients/${id}`);
  }

  create(client: Partial<Client>): Observable<Client> {
    return this.api.post<Client>('/clients', client);
  }

  update(id: number, client: Partial<Client>): Observable<Client> {
    return this.api.put<Client>(`/clients/${id}`, client);
  }

  delete(id: number): Observable<void> {
    return this.api.delete<void>(`/clients/${id}`);
  }
}
