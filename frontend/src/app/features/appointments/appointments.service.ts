import { Injectable, inject } from '@angular/core';
import { Observable } from 'rxjs';
import { ApiService } from '../../core/services/api.service';
import { Appointment } from '../../core/models';

@Injectable({ providedIn: 'root' })
export class AppointmentsService {
  private readonly api = inject(ApiService);

  getAll(): Observable<Appointment[]> {
    return this.api.get<Appointment[]>('/appointments');
  }

  getById(id: string): Observable<Appointment> {
    return this.api.get<Appointment>(`/appointments/${id}`);
  }

  create(appointment: Partial<Appointment>): Observable<Appointment> {
    return this.api.post<Appointment>('/appointments', appointment);
  }

  update(id: string, appointment: Partial<Appointment>): Observable<Appointment> {
    return this.api.put<Appointment>(`/appointments/${id}`, appointment);
  }

  delete(id: string): Observable<void> {
    return this.api.delete<void>(`/appointments/${id}`);
  }
}
