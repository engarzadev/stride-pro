import { Component, computed, inject, OnInit, signal } from '@angular/core';
import { Router } from '@angular/router';
import { MatCardModule } from '@angular/material/card';
import { Appointment } from '../../../core/models';
import { ConfirmDialogService } from '../../../shared/components/confirm-dialog/confirm-dialog.component';
import {
  DataTableComponent,
  MobileCardConfig,
  TableAction,
  TableColumn,
} from '../../../shared/components/data-table/data-table.component';
import { LoadingSpinnerComponent } from '../../../shared/components/loading-spinner/loading-spinner.component';
import { PageHeaderComponent } from '../../../shared/components/page-header/page-header.component';
import { ToastService } from '../../../shared/components/toast/toast.service';
import { AppointmentsService } from '../appointments.service';

@Component({
  selector: 'app-appointment-list',
  standalone: true,
  imports: [PageHeaderComponent, DataTableComponent, LoadingSpinnerComponent, MatCardModule],
  templateUrl: './appointment-list.component.html',
  styleUrls: ['./appointment-list.component.scss'],
})
export class AppointmentListComponent implements OnInit {
  private readonly appointmentsService = inject(AppointmentsService);
  private readonly router = inject(Router);
  private readonly confirmDialog = inject(ConfirmDialogService);
  private readonly toast = inject(ToastService);

  readonly loading = signal(true);
  readonly appointments = signal<Appointment[]>([]);

  readonly tableAppointments = computed(() =>
    this.appointments().map((a) => ({
      ...a,
      time: a.time
        ? (() => {
            const [h, m] = a.time!.split(':');
            const hour = parseInt(h, 10);
            return `${hour % 12 || 12}:${m} ${hour < 12 ? 'AM' : 'PM'}`;
          })()
        : '',
      clientName: a.client ? `${a.client.firstName} ${a.client.lastName}` : '',
      horseName: a.horse?.name ?? '',
    })),
  );

  readonly columns: TableColumn[] = [
    { key: 'date', label: 'Date', sortable: true, type: 'date' },
    { key: 'time', label: 'Time' },
    { key: 'clientName', label: 'Client', sortable: true },
    { key: 'horseName', label: 'Horse', sortable: true },
    { key: 'type', label: 'Type', capitalize: true },
    {
      key: 'status',
      label: 'Status',
      type: 'badge',
      capitalize: true,
      badgeMap: {
        scheduled: 'success',
        confirmed: 'success',
        completed: 'secondary',
        cancelled: 'danger',
        'no-show': 'warning',
      },
    },
  ];

  readonly actions: TableAction[] = [
    { label: 'Edit', action: 'edit', class: 'btn-outline' },
    { label: 'Delete', action: 'delete', class: 'btn-danger' },
  ];

  readonly mobileCard: MobileCardConfig = { titleKey: 'date', subtitleKey: 'horseName' };

  ngOnInit(): void {
    this.loadAppointments();
  }

  loadAppointments(): void {
    this.appointmentsService.getAll().subscribe({
      next: (appointments) => {
        this.appointments.set(appointments);
        this.loading.set(false);
      },
      error: () => this.loading.set(false),
    });
  }

  onAdd(): void {
    this.router.navigate(['/appointments/new']);
  }

  onRowClick(row: Record<string, unknown>): void {
    this.router.navigate(['/appointments', row['id']]);
  }

  async onAction(event: {
    action: string;
    row: Record<string, unknown>;
  }): Promise<void> {
    if (event.action === 'edit') {
      this.router.navigate(['/appointments', event.row['id'], 'edit']);
    } else if (event.action === 'delete') {
      const confirmed = await this.confirmDialog.confirm({
        title: 'Delete Appointment',
        message:
          'Are you sure you want to delete this appointment? This action cannot be undone.',
        confirmText: 'Delete',
        confirmClass: 'btn-danger',
      });
      if (confirmed) {
        this.appointmentsService.delete(event.row['id'] as string).subscribe({
          next: () => {
            this.toast.success('Appointment deleted successfully');
            this.loadAppointments();
          },
        });
      }
    }
  }
}
