import { Component, inject, OnInit, signal } from '@angular/core';
import { ActivatedRoute, Router, RouterLink } from '@angular/router';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { SessionsService } from '../sessions.service';
import { AppointmentsService } from '../../appointments/appointments.service';
import { Appointment } from '../../../core/models';
import { ToastService } from '../../../shared/components/toast/toast.service';
import { LoadingSpinnerComponent } from '../../../shared/components/loading-spinner/loading-spinner.component';
import { DateFormatPipe } from '../../../shared/pipes/date-format.pipe';
import { QuickCreateAppointmentService } from '../../../shared/components/quick-create/quick-create-appointment.component';

@Component({
  selector: 'app-session-form',
  standalone: true,
  imports: [ReactiveFormsModule, RouterLink, LoadingSpinnerComponent, DateFormatPipe],
  templateUrl: './session-form.component.html',
  styleUrls: ['./session-form.component.scss'],
})
export class SessionFormComponent implements OnInit {
  private readonly fb = inject(FormBuilder);
  private readonly route = inject(ActivatedRoute);
  private readonly router = inject(Router);
  private readonly sessionsService = inject(SessionsService);
  private readonly appointmentsService = inject(AppointmentsService);
  private readonly toast = inject(ToastService);
  private readonly quickCreateAppointment = inject(QuickCreateAppointmentService);

  readonly loading = signal(false);
  readonly saving = signal(false);
  readonly isEdit = signal(false);
  readonly appointments = signal<Appointment[]>([]);
  private sessionId = '';

  readonly bodyZoneOptions = [
    'Head', 'Neck', 'Withers', 'Back', 'Loin', 'Croup',
    'Shoulder', 'Foreleg', 'Hindleg', 'Hoof', 'Barrel',
    'Chest', 'Abdomen', 'Hip', 'Stifle', 'Hock', 'Fetlock',
    'Poll', 'TMJ', 'Pelvis', 'Sacrum',
  ];

  readonly selectedZones = signal<Set<string>>(new Set());

  readonly form = this.fb.nonNullable.group({
    appointmentId: ['', [Validators.required]],
    type: ['', [Validators.required]],
    notes: [''],
    findings: [''],
    recommendations: [''],
  });

  ngOnInit(): void {
    this.appointmentsService.getAll().subscribe((a) => this.appointments.set(a));

    const id = this.route.snapshot.paramMap.get('id');
    if (id) {
      this.isEdit.set(true);
      this.sessionId = id;
      this.loading.set(true);
      this.sessionsService.getById(this.sessionId).subscribe({
        next: (session) => {
          this.form.patchValue({
            appointmentId: session.appointmentId,
            type: session.type,
            notes: session.notes,
            findings: session.findings,
            recommendations: session.recommendations,
          });
          if (session.bodyZones) {
            this.selectedZones.set(new Set(session.bodyZones));
          }
          this.loading.set(false);
        },
        error: () => {
          this.loading.set(false);
          this.router.navigate(['/sessions']);
        },
      });
    }
  }

  async openCreateAppointment(): Promise<void> {
    const appointment = await this.quickCreateAppointment.open();
    if (appointment) {
      this.appointments.update((a) => [...a, appointment]);
      this.form.controls.appointmentId.setValue(appointment.id);
    }
  }

  toggleZone(zone: string): void {
    this.selectedZones.update((zones) => {
      const next = new Set(zones);
      if (next.has(zone)) {
        next.delete(zone);
      } else {
        next.add(zone);
      }
      return next;
    });
  }

  onSubmit(): void {
    if (this.form.invalid) {
      this.form.markAllAsTouched();
      return;
    }

    this.saving.set(true);
    const formData = this.form.getRawValue();
    const data = {
      ...formData,
      bodyZones: Array.from(this.selectedZones()),
    };

    const request$ = this.isEdit()
      ? this.sessionsService.update(this.sessionId, data)
      : this.sessionsService.create(data);

    request$.subscribe({
      next: () => {
        this.toast.success(this.isEdit() ? 'Session updated successfully' : 'Session created successfully');
        this.router.navigate(['/sessions']);
      },
      error: () => this.saving.set(false),
    });
  }
}
