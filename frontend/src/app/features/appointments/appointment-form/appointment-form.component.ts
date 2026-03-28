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
import { QuickCreateClientService } from '../../../shared/components/quick-create/quick-create-client.component';
import { QuickCreateHorseService } from '../../../shared/components/quick-create/quick-create-horse.component';
import { QuickCreateBarnService } from '../../../shared/components/quick-create/quick-create-barn.component';

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
  private readonly quickCreateClient = inject(QuickCreateClientService);
  private readonly quickCreateHorse = inject(QuickCreateHorseService);
  private readonly quickCreateBarn = inject(QuickCreateBarnService);

  readonly loading = signal(false);
  readonly saving = signal(false);
  readonly isEdit = signal(false);
  readonly clients = signal<Client[]>([]);
  readonly allHorses = signal<Horse[]>([]);
  readonly barns = signal<Barn[]>([]);
  private appointmentId = '';

  readonly form = this.fb.nonNullable.group({
    clientId: ['', [Validators.required]],
    horseId: ['', [Validators.required]],
    barnId: [null as string | null],
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
    return this.allHorses().filter((h) => h.clientId === clientId);
  }

  ngOnInit(): void {
    this.clientsService.getAll().subscribe((c) => this.clients.set(c));
    this.horsesService.getAll().subscribe((h) => this.allHorses.set(h));
    this.barnsService.getAll().subscribe((b) => this.barns.set(b));

    this.form.controls.horseId.disable();
    this.form.controls.barnId.disable();

    this.form.controls.clientId.valueChanges.subscribe((clientId) => {
      this.form.controls.horseId.setValue('');
      this.form.controls.barnId.setValue(null);
      this.form.controls.barnId.disable();
      if (clientId) {
        this.form.controls.horseId.enable();
      } else {
        this.form.controls.horseId.disable();
      }
    });

    this.form.controls.horseId.valueChanges.subscribe((horseId) => {
      const horse = this.allHorses().find((h) => h.id === horseId);
      if (horseId) {
        this.form.controls.barnId.enable();
        this.form.controls.barnId.setValue(horse?.barnId ?? null);
      } else {
        this.form.controls.barnId.setValue(null);
        this.form.controls.barnId.disable();
      }
    });

    const id = this.route.snapshot.paramMap.get('id');
    if (id) {
      this.isEdit.set(true);
      this.appointmentId = id;
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

  async openCreateClient(): Promise<void> {
    const client = await this.quickCreateClient.open();
    if (client) {
      this.clients.update((c) => [...c, client]);
      this.form.controls.clientId.setValue(client.id);
    }
  }

  async openCreateHorse(): Promise<void> {
    const selectedClientId = this.form.controls.clientId.value;
    const horse = await this.quickCreateHorse.open({ clientId: selectedClientId || undefined });
    if (horse) {
      this.allHorses.update((h) => [...h, horse]);
      this.form.controls.horseId.setValue(horse.id);
      if (horse.barnId) {
        this.form.controls.barnId.setValue(horse.barnId);
      }
    }
  }

  async openCreateBarn(): Promise<void> {
    const barn = await this.quickCreateBarn.open();
    if (barn) {
      this.barns.update((b) => [...b, barn]);
      this.form.controls.barnId.setValue(barn.id);
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
