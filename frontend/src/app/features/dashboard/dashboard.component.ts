import { Component, inject, OnInit, signal } from '@angular/core';
import { Router } from '@angular/router';
import { forkJoin } from 'rxjs';
import { Appointment, Client, Horse, Invoice } from '../../core/models';
import { ApiService } from '../../core/services/api.service';
import { LoadingSpinnerComponent } from '../../shared/components/loading-spinner/loading-spinner.component';
import { DateFormatPipe } from '../../shared/pipes/date-format.pipe';

@Component({
  selector: 'app-dashboard',
  standalone: true,
  imports: [LoadingSpinnerComponent, DateFormatPipe],
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.scss'],
})
export class DashboardComponent implements OnInit {
  private readonly api = inject(ApiService);
  private readonly router = inject(Router);

  readonly loading = signal(true);
  readonly clients = signal<Client[]>([]);
  readonly horses = signal<Horse[]>([]);
  readonly appointments = signal<Appointment[]>([]);
  readonly invoices = signal<Invoice[]>([]);

  ngOnInit(): void {
    forkJoin({
      clients: this.api.get<Client[]>('/clients'),
      horses: this.api.get<Horse[]>('/horses'),
      appointments: this.api.get<Appointment[]>('/appointments'),
      invoices: this.api.get<Invoice[]>('/invoices'),
    }).subscribe({
      next: (data) => {
        this.clients.set(data.clients || []);
        this.horses.set(data.horses || []);
        this.appointments.set(data.appointments || []);
        this.invoices.set(data.invoices || []);
        this.loading.set(false);
      },
      error: () => {
        this.loading.set(false);
      },
    });
  }

  get upcomingAppointments(): Appointment[] {
    const now = new Date().toISOString();
    return this.appointments()
      .filter((a) => a.date >= now || a.status === 'scheduled')
      .slice(0, 5);
  }

  get pendingInvoices(): Invoice[] {
    return this.invoices().filter(
      (i) => i.status === 'pending' || i.status === 'overdue',
    );
  }

  navigateTo(path: string): void {
    this.router.navigate([path]);
  }
}
