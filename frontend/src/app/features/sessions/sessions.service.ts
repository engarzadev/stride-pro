import { Injectable, inject } from '@angular/core';
import { Observable } from 'rxjs';
import { ApiService } from '../../core/services/api.service';
import { Session } from '../../core/models';

@Injectable({ providedIn: 'root' })
export class SessionsService {
  private readonly api = inject(ApiService);

  getAll(): Observable<Session[]> {
    return this.api.get<Session[]>('/sessions');
  }

  getById(id: string): Observable<Session> {
    return this.api.get<Session>(`/sessions/${id}`);
  }

  create(session: Partial<Session>): Observable<Session> {
    return this.api.post<Session>('/sessions', session);
  }

  update(id: string, session: Partial<Session>): Observable<Session> {
    return this.api.put<Session>(`/sessions/${id}`, session);
  }

  delete(id: string): Observable<void> {
    return this.api.delete<void>(`/sessions/${id}`);
  }
}
