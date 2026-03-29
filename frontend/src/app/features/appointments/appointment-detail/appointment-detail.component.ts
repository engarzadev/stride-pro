import { Component, inject, OnInit, signal } from '@angular/core';
import { TitleCasePipe } from '@angular/common';
import { ActivatedRoute, Router, RouterLink } from '@angular/router';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { AppointmentsService } from '../appointments.service';
import { Appointment } from '../../../core/models';
import { LoadingSpinnerComponent } from '../../../shared/components/loading-spinner/loading-spinner.component';
import { ConfirmDialogService } from '../../../shared/components/confirm-dialog/confirm-dialog.component';
import { ToastService } from '../../../shared/components/toast/toast.service';
import { DateFormatPipe } from '../../../shared/pipes/date-format.pipe';

@Component({
  selector: 'app-appointment-detail',
  standalone: true,
  imports: [RouterLink, LoadingSpinnerComponent, DateFormatPipe, TitleCasePipe, MatCardModule, MatButtonModule, MatIconModule],
  templateUrl: './appointment-detail.component.html',
  styleUrls: ['./appointment-detail.component.scss'],
})
export class AppointmentDetailComponent implements OnInit {
  private readonly route = inject(ActivatedRoute);
  private readonly router = inject(Router);
  private readonly appointmentsService = inject(AppointmentsService);
  private readonly confirmDialog = inject(ConfirmDialogService);
  private readonly toast = inject(ToastService);

  readonly loading = signal(true);
  readonly appointment = signal<Appointment | null>(null);

  ngOnInit(): void {
    const id = this.route.snapshot.paramMap.get('id')!;
    this.appointmentsService.getById(id).subscribe({
      next: (appointment) => {
        this.appointment.set(appointment);
        this.loading.set(false);
      },
      error: () => {
        this.loading.set(false);
        this.router.navigate(['/appointments']);
      },
    });
  }

  getStatusClass(status: string): string {
    const map: Record<string, string> = {
      scheduled: 'badge-success',
      confirmed: 'badge-success',
      completed: 'badge-secondary',
      cancelled: 'badge-danger',
      'no-show': 'badge-warning',
    };
    return map[status] || 'badge-secondary';
  }

  async onDelete(): Promise<void> {
    const a = this.appointment();
    if (!a) return;

    const confirmed = await this.confirmDialog.confirm({
      title: 'Delete Appointment',
      message: 'Are you sure you want to delete this appointment?',
      confirmText: 'Delete',
      confirmClass: 'btn-danger',
    });

    if (confirmed) {
      this.appointmentsService.delete(a.id).subscribe({
        next: () => {
          this.toast.success('Appointment deleted successfully');
          this.router.navigate(['/appointments']);
        },
      });
    }
  }
}
