import { Injectable, inject } from '@angular/core';
import { Observable } from 'rxjs';
import { ApiService } from '../../core/services/api.service';
import { Barn } from '../../core/models';

@Injectable({ providedIn: 'root' })
export class BarnsService {
  private readonly api = inject(ApiService);

  getAll(): Observable<Barn[]> {
    return this.api.get<Barn[]>('/barns');
  }

  getById(id: string): Observable<Barn> {
    return this.api.get<Barn>(`/barns/${id}`);
  }

  create(barn: Partial<Barn>): Observable<Barn> {
    return this.api.post<Barn>('/barns', barn);
  }

  update(id: string, barn: Partial<Barn>): Observable<Barn> {
    return this.api.put<Barn>(`/barns/${id}`, barn);
  }

  delete(id: string): Observable<void> {
    return this.api.delete<void>(`/barns/${id}`);
  }
}
