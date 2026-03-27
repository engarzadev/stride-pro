import { Component, inject, OnInit, signal } from '@angular/core';
import { ActivatedRoute, Router, RouterLink } from '@angular/router';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { AppointmentsService } from '../appointments.service';
import { ClientsService } from '../../clients/clients.service';
import { HorsesService } from '../../horses/horses.service';
import { BarnsService } from '../../barns/barns.service';
import { Client, Horse, Barn } from '../../../core/models';
import { ToastService } from '../../../shared/components/toast/toast.service';
import { LoadingSpinnerComponent } from '../../../shared/components/loading-spinner/loading-spinner.component';

@Component({
  selector: 'app-appointment-form',
  standalone: true,
  imports: [ReactiveFormsModule, RouterLink, LoadingSpinnerComponent],
  templateUrl: './appointment-form.component.html',
  styleUrls: ['./appointment-form.component.scss'],
})
export class AppointmentFormComponent implements OnInit {
  private readonly fb = inject(FormBuilder);
  private readonly route = inject(ActivatedRoute);
  private readonly router = inject(Router);
  private readonly appointmentsService = inject(AppointmentsService);
  private readonly clientsService = inject(ClientsService);
  private readonly horsesService = inject(HorsesService);
  private readonly barnsService = inject(BarnsService);
  private readonly toast = inject(ToastService);

  readonly loading = signal(false);
  readonly saving = signal(false);
  readonly isEdit = signal(false);
  readonly clients = signal<Client[]>([]);
  readonly allHorses = signal<Horse[]>([]);
  readonly barns = signal<Barn[]>([]);
  private appointmentId = 0;

  readonly form = this.fb.nonNullable.group({
    clientId: [0, [Validators.required, Validators.min(1)]],
    horseId: [0, [Validators.required, Validators.min(1)]],
    barnId: [0],
    date: ['', [Validators.required]],
    time: [''],
    duration: [60],
    type: ['', [Validators.required]],
    status: ['scheduled'],
    notes: [''],
  });

  get filteredHorses(): Horse[] {
    const clientId = this.form.controls.clientId.value;
    if (!clientId) return this.allHorses();
    return this.allHorses().filter((h) => h.clientId === Number(clientId));
  }

  ngOnInit(): void {
    this.clientsService.getAll().subscribe((c) => this.clients.set(c));
    this.horsesService.getAll().subscribe((h) => this.allHorses.set(h));
    this.barnsService.getAll().subscribe((b) => this.barns.set(b));

    this.form.controls.clientId.valueChanges.subscribe(() => {
      this.form.controls.horseId.setValue(0);
    });

    this.form.controls.horseId.valueChanges.subscribe((horseId) => {
      const horse = this.allHorses().find((h) => h.id === Number(horseId));
      if (horse?.barnId) {
        this.form.controls.barnId.setValue(horse.barnId);
      }
    });

    const id = this.route.snapshot.paramMap.get('id');
    if (id) {
      this.isEdit.set(true);
      this.appointmentId = Number(id);
      this.loading.set(true);
      this.appointmentsService.getById(this.appointmentId).subscribe({
        next: (apt) => {
          this.form.patchValue({
            clientId: apt.clientId,
            horseId: apt.horseId,
            barnId: apt.barnId,
            date: apt.date?.substring(0, 10),
            time: apt.time,
            duration: apt.duration,
            type: apt.type,
            status: apt.status,
            notes: apt.notes,
          });
          this.loading.set(false);
        },
        error: () => {
          this.loading.set(false);
          this.router.navigate(['/appointments']);
        },
      });
    }
  }

  onSubmit(): void {
    if (this.form.invalid) {
      this.form.markAllAsTouched();
      return;
    }

    this.saving.set(true);
    const data = this.form.getRawValue();

    const request$ = this.isEdit()
      ? this.appointmentsService.update(this.appointmentId, data)
      : this.appointmentsService.create(data);

    request$.subscribe({
      next: () => {
        this.toast.success(this.isEdit() ? 'Appointment updated successfully' : 'Appointment created successfully');
        this.router.navigate(['/appointments']);
      },
      error: () => this.saving.set(false),
    });
  }
}
