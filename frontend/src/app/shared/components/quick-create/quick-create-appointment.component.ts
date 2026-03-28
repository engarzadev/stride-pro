import { Component, Injectable, OnInit, inject, signal } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { Client, Horse, Barn, Appointment } from '../../../core/models';
import { AppointmentsService } from '../../../features/appointments/appointments.service';
import { ClientsService } from '../../../features/clients/clients.service';
import { HorsesService } from '../../../features/horses/horses.service';
import { BarnsService } from '../../../features/barns/barns.service';

@Injectable({ providedIn: 'root' })
export class QuickCreateAppointmentService {
  readonly visible = signal(false);
  private resolveFn?: (result: Appointment | null) => void;

  open(): Promise<Appointment | null> {
    this.visible.set(true);
    return new Promise<Appointment | null>((resolve) => {
      this.resolveFn = resolve;
    });
  }

  complete(appointment: Appointment): void {
    this.visible.set(false);
    this.resolveFn?.(appointment);
  }

  cancel(): void {
    this.visible.set(false);
    this.resolveFn?.(null);
  }
}

@Component({
  selector: 'app-quick-create-appointment',
  standalone: true,
  imports: [ReactiveFormsModule],
  templateUrl: './quick-create-appointment.component.html',
  styleUrls: ['./quick-create-modal.scss'],
})
export class QuickCreateAppointmentComponent implements OnInit {
  readonly service = inject(QuickCreateAppointmentService);
  private readonly appointmentsService = inject(AppointmentsService);
  private readonly clientsService = inject(ClientsService);
  private readonly horsesService = inject(HorsesService);
  private readonly barnsService = inject(BarnsService);
  private readonly fb = inject(FormBuilder);

  readonly saving = signal(false);
  readonly clients = signal<Client[]>([]);
  readonly allHorses = signal<Horse[]>([]);
  readonly barns = signal<Barn[]>([]);

  readonly form = this.fb.nonNullable.group({
    clientId: ['', [Validators.required]],
    horseId: ['', [Validators.required]],
    barnId: [null as string | null],
    date: ['', [Validators.required]],
    time: [''],
    duration: [60],
    type: ['', [Validators.required]],
    notes: [''],
  });

  get filteredHorses(): Horse[] {
    const clientId = this.form.controls.clientId.value;
    if (!clientId) return this.allHorses();
    return this.allHorses().filter((h) => h.clientId === clientId);
  }

  ngOnInit(): void {
    this.clientsService.getAll().subscribe((c) => this.clients.set(c));
    this.horsesService.getAll().subscribe((h) => this.allHorses.set(h));
    this.barnsService.getAll().subscribe((b) => this.barns.set(b));

    this.form.controls.clientId.valueChanges.subscribe(() => {
      this.form.controls.horseId.setValue('');
    });

    this.form.controls.horseId.valueChanges.subscribe((horseId) => {
      const horse = this.allHorses().find((h) => h.id === horseId);
      if (horse?.barnId) {
        this.form.controls.barnId.setValue(horse.barnId);
      }
    });
  }

  onSubmit(): void {
    if (this.form.invalid) {
      this.form.markAllAsTouched();
      return;
    }
    this.saving.set(true);
    const data = this.form.getRawValue();
    this.appointmentsService.create(data).subscribe({
      next: (appointment) => {
        this.saving.set(false);
        this.form.reset({ clientId: '', horseId: '', barnId: null, date: '', time: '', duration: 60, type: '', notes: '' });
        this.service.complete(appointment);
      },
      error: () => this.saving.set(false),
    });
  }

  onCancel(): void {
    this.form.reset({ clientId: '', horseId: '', barnId: null, date: '', time: '', duration: 60, type: '', notes: '' });
    this.service.cancel();
  }
}
