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

  getById(id: number): Observable<Appointment> {
    return this.api.get<Appointment>(`/appointments/${id}`);
  }

  create(appointment: Partial<Appointment>): Observable<Appointment> {
    return this.api.post<Appointment>('/appointments', appointment);
  }

  update(id: number, appointment: Partial<Appointment>): Observable<Appointment> {
    return this.api.put<Appointment>(`/appointments/${id}`, appointment);
  }

  delete(id: number): Observable<void> {
    return this.api.delete<void>(`/appointments/${id}`);
  }
}
