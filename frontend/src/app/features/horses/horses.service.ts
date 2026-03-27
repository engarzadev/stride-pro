import { Injectable, inject } from '@angular/core';
import { Observable } from 'rxjs';
import { ApiService } from '../../core/services/api.service';
import { Horse } from '../../core/models';

@Injectable({ providedIn: 'root' })
export class HorsesService {
  private readonly api = inject(ApiService);

  getAll(): Observable<Horse[]> {
    return this.api.get<Horse[]>('/horses');
  }

  getById(id: number): Observable<Horse> {
    return this.api.get<Horse>(`/horses/${id}`);
  }

  create(horse: Partial<Horse>): Observable<Horse> {
    return this.api.post<Horse>('/horses', horse);
  }

  update(id: number, horse: Partial<Horse>): Observable<Horse> {
    return this.api.put<Horse>(`/horses/${id}`, horse);
  }

  delete(id: number): Observable<void> {
    return this.api.delete<void>(`/horses/${id}`);
  }
}
