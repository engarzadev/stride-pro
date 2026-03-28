import { Component, Injectable, OnInit, inject, signal } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import {
  MatDialog,
  MatDialogModule,
  MatDialogRef,
} from '@angular/material/dialog';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';
import { firstValueFrom } from 'rxjs';
import { Appointment, Barn, Client, Horse } from '../../../core/models';
import { AppointmentsService } from '../../../features/appointments/appointments.service';
import { BarnsService } from '../../../features/barns/barns.service';
import { ClientsService } from '../../../features/clients/clients.service';
import { HorsesService } from '../../../features/horses/horses.service';

@Injectable({ providedIn: 'root' })
export class QuickCreateAppointmentService {
  private readonly dialog = inject(MatDialog);

  async open(): Promise<Appointment | null> {
    const ref = this.dialog.open(QuickCreateAppointmentComponent, {
      width: '650px',
    });
    return (await firstValueFrom(ref.afterClosed())) ?? null;
  }
}

@Component({
  selector: 'app-quick-create-appointment',
  standalone: true,
  imports: [
    ReactiveFormsModule,
    MatDialogModule,
    MatFormFieldModule,
    MatInputModule,
    MatSelectModule,
    MatButtonModule,
  ],
  template: `
    <h2 mat-dialog-title>New Appointment</h2>
    <mat-dialog-content>
      <form [formGroup]="form" id="qc-apt-form" (ngSubmit)="onSubmit()">
        <div class="form-row">
          <mat-form-field appearance="outline">
            <mat-label>Client</mat-label>
            <mat-select formControlName="clientId">
              <mat-option [value]="0">Select a client</mat-option>
              @for (client of clients(); track client.id) {
                <mat-option [value]="client.id"
                  >{{ client.firstName }} {{ client.lastName }}</mat-option
                >
              }
            </mat-select>
            @if (form.controls.clientId.errors) {
              <mat-error>Please select a client.</mat-error>
            }
          </mat-form-field>
          <mat-form-field appearance="outline">
            <mat-label>Horse</mat-label>
            <mat-select formControlName="horseId">
              <mat-option [value]="0">Select a horse</mat-option>
              @for (horse of filteredHorses; track horse.id) {
                <mat-option [value]="horse.id">{{ horse.name }}</mat-option>
              }
            </mat-select>
            @if (form.controls.horseId.errors) {
              <mat-error>Please select a horse.</mat-error>
            }
          </mat-form-field>
        </div>
        <mat-form-field appearance="outline">
          <mat-label>Barn</mat-label>
          <mat-select formControlName="barnId">
            <mat-option [value]="0">No barn</mat-option>
            @for (barn of barns(); track barn.id) {
              <mat-option [value]="barn.id">{{ barn.name }}</mat-option>
            }
          </mat-select>
        </mat-form-field>
        <div class="form-row">
          <mat-form-field appearance="outline">
            <mat-label>Date</mat-label>
            <input matInput type="date" formControlName="date" />
            @if (form.controls.date.errors?.['required']) {
              <mat-error>Date is required.</mat-error>
            }
          </mat-form-field>
          <mat-form-field appearance="outline">
            <mat-label>Time</mat-label>
            <input matInput type="time" formControlName="time" />
          </mat-form-field>
        </div>
        <div class="form-row">
          <mat-form-field appearance="outline">
            <mat-label>Duration (minutes)</mat-label>
            <input
              matInput
              type="number"
              formControlName="duration"
              min="15"
              step="15"
            />
          </mat-form-field>
          <mat-form-field appearance="outline">
            <mat-label>Type</mat-label>
            <mat-select formControlName="type">
              <mat-option value="">Select type</mat-option>
              <mat-option value="checkup">Checkup</mat-option>
              <mat-option value="treatment">Treatment</mat-option>
              <mat-option value="massage">Massage</mat-option>
              <mat-option value="chiropractic">Chiropractic</mat-option>
              <mat-option value="acupuncture">Acupuncture</mat-option>
              <mat-option value="rehabilitation">Rehabilitation</mat-option>
              <mat-option value="evaluation">Evaluation</mat-option>
              <mat-option value="other">Other</mat-option>
            </mat-select>
            @if (form.controls.type.errors?.['required']) {
              <mat-error>Type is required.</mat-error>
            }
          </mat-form-field>
        </div>
        <mat-form-field appearance="outline">
          <mat-label>Notes</mat-label>
          <textarea matInput formControlName="notes" rows="2"></textarea>
        </mat-form-field>
      </form>
    </mat-dialog-content>
    <mat-dialog-actions align="end">
      <button mat-stroked-button (click)="dialogRef.close(null)">Cancel</button>
      <button
        mat-raised-button
        color="primary"
        form="qc-apt-form"
        type="submit"
        [disabled]="saving()"
      >
        @if (saving()) {
          Saving...
        } @else {
          Create Appointment
        }
      </button>
    </mat-dialog-actions>
  `,
})
export class QuickCreateAppointmentComponent implements OnInit {
  readonly dialogRef = inject(MatDialogRef<QuickCreateAppointmentComponent>);
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
    this.appointmentsService.create(this.form.getRawValue()).subscribe({
      next: (appointment) => this.dialogRef.close(appointment),
      error: () => this.saving.set(false),
    });
  }
}
