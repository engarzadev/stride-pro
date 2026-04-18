import { Injectable, inject } from '@angular/core';
import { Observable } from 'rxjs';
import { ApiService } from '../../core/services/api.service';
import { CareLog, Horse, Reminder } from '../../core/models';

@Injectable({ providedIn: 'root' })
export class HorsesService {
  private readonly api = inject(ApiService);

  getAll(): Observable<Horse[]> {
    return this.api.get<Horse[]>('/horses');
  }

  getById(id: string): Observable<Horse> {
    return this.api.get<Horse>(`/horses/${id}`);
  }

  create(horse: Partial<Horse>): Observable<Horse> {
    return this.api.post<Horse>('/horses', horse);
  }

  update(id: string, horse: Partial<Horse>): Observable<Horse> {
    return this.api.put<Horse>(`/horses/${id}`, horse);
  }

  delete(id: string): Observable<void> {
    return this.api.delete<void>(`/horses/${id}`);
  }

  getCareLogs(horseId: string): Observable<CareLog[]> {
    return this.api.get<CareLog[]>(`/horses/${horseId}/care-logs`);
  }

  createCareLog(horseId: string, entry: Partial<CareLog>): Observable<CareLog> {
    return this.api.post<CareLog>(`/horses/${horseId}/care-logs`, entry);
  }

  updateCareLog(id: string, entry: Partial<CareLog>): Observable<CareLog> {
    return this.api.put<CareLog>(`/care-logs/${id}`, entry);
  }

  deleteCareLog(id: string): Observable<void> {
    return this.api.delete<void>(`/care-logs/${id}`);
  }

  getReminders(horseId: string): Observable<Reminder[]> {
    return this.api.get<Reminder[]>(`/horses/${horseId}/reminders`);
  }

  createReminder(horseId: string, reminder: Partial<Reminder> & { isComplete?: boolean }): Observable<Reminder> {
    return this.api.post<Reminder>(`/horses/${horseId}/reminders`, reminder);
  }

  putReminder(id: string, reminder: Partial<Reminder>): Observable<Reminder> {
    return this.api.put<Reminder>(`/reminders/${id}`, reminder);
  }

  patchReminder(id: string, patch: { isComplete: boolean }): Observable<Reminder> {
    return this.api.patch<Reminder>(`/reminders/${id}`, patch);
  }

  deleteReminder(id: string): Observable<void> {
    return this.api.delete<void>(`/reminders/${id}`);
  }
}
