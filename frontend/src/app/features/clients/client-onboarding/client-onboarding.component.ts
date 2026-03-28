import { DatePipe, TitleCasePipe } from '@angular/common';
import { Component, inject, OnInit, signal } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { Router, RouterLink } from '@angular/router';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';
import { Appointment, Barn, Client, Horse } from '../../../core/models';
import { BreedAutocompleteComponent } from '../../../shared/components/breed-autocomplete/breed-autocomplete.component';
import { QuickCreateBarnService } from '../../../shared/components/quick-create/quick-create-barn.component';
import { ToastService } from '../../../shared/components/toast/toast.service';
import { AppointmentsService } from '../../appointments/appointments.service';
import { BarnsService } from '../../barns/barns.service';
import { HorsesService } from '../../horses/horses.service';
import { ClientsService } from '../clients.service';

@Component({
  selector: 'app-client-onboarding',
  standalone: true,
  imports: [
    ReactiveFormsModule,
    RouterLink,
    BreedAutocompleteComponent,
    TitleCasePipe,
    DatePipe,
    MatButtonModule,
    MatCardModule,
    MatFormFieldModule,
    MatIconModule,
    MatInputModule,
    MatSelectModule,
  ],
  templateUrl: './client-onboarding.component.html',
  styleUrls: ['./client-onboarding.component.scss'],
})
export class ClientOnboardingComponent implements OnInit {
  private readonly fb = inject(FormBuilder);
  readonly router = inject(Router);
  private readonly clientsService = inject(ClientsService);
  private readonly horsesService = inject(HorsesService);
  private readonly barnsService = inject(BarnsService);
  private readonly appointmentsService = inject(AppointmentsService);
  private readonly toast = inject(ToastService);
  private readonly quickCreateBarn = inject(QuickCreateBarnService);

  readonly step = signal(1);
  readonly saving = signal(false);
  readonly barns = signal<Barn[]>([]);

  readonly createdClient = signal<Client | null>(null);
  readonly createdHorse = signal<Horse | null>(null);
  readonly createdAppointment = signal<Appointment | null>(null);

  readonly clientForm = this.fb.nonNullable.group({
    firstName: ['', [Validators.required]],
    lastName: ['', [Validators.required]],
    email: ['', [Validators.email]],
    phone: [''],
    address: [''],
    notes: [''],
  });

  readonly horseForm = this.fb.nonNullable.group({
    name: ['', [Validators.required]],
    breed: [''],
    age: [null as number | null],
    gender: [''],
    color: [''],
    weight: [null as number | null],
    barnId: [null as string | null],
    notes: [''],
  });

  readonly appointmentForm = this.fb.nonNullable.group({
    date: ['', [Validators.required]],
    time: [''],
    duration: [60],
    type: ['', [Validators.required]],
    notes: [''],
  });

  readonly appointmentTypes = [
    'checkup',
    'treatment',
    'massage',
    'chiropractic',
    'acupuncture',
    'rehabilitation',
    'evaluation',
    'other',
  ];

  ngOnInit(): void {
    this.barnsService.getAll().subscribe((b) => this.barns.set(b));
  }

  saveClient(): void {
    if (this.clientForm.invalid) {
      this.clientForm.markAllAsTouched();
      return;
    }
    this.saving.set(true);
    this.clientsService.create(this.clientForm.getRawValue()).subscribe({
      next: (client) => {
        this.createdClient.set(client);
        this.saving.set(false);
        this.step.set(2);
      },
      error: () => this.saving.set(false),
    });
  }

  saveHorse(): void {
    if (this.horseForm.invalid) {
      this.horseForm.markAllAsTouched();
      return;
    }
    this.saving.set(true);
    const raw = this.horseForm.getRawValue();
    const data = {
      ...raw,
      clientId: this.createdClient()!.id,
      age: raw.age ?? undefined,
      weight: raw.weight ?? undefined,
      barnId: raw.barnId ?? undefined,
    };
    this.horsesService.create(data).subscribe({
      next: (horse) => {
        this.createdHorse.set(horse);
        this.saving.set(false);
        this.step.set(3);
      },
      error: () => this.saving.set(false),
    });
  }

  skipHorse(): void {
    this.step.set(4);
  }

  saveAppointment(): void {
    if (this.appointmentForm.invalid) {
      this.appointmentForm.markAllAsTouched();
      return;
    }
    this.saving.set(true);
    const horse = this.createdHorse()!;
    const data = {
      ...this.appointmentForm.getRawValue(),
      clientId: this.createdClient()!.id,
      horseId: horse.id,
      barnId: horse.barnId ?? null,
      status: 'scheduled',
    };
    this.appointmentsService.create(data).subscribe({
      next: (apt) => {
        this.createdAppointment.set(apt);
        this.saving.set(false);
        this.step.set(4);
      },
      error: () => this.saving.set(false),
    });
  }

  skipAppointment(): void {
    this.step.set(4);
  }

  async addNewBarn(): Promise<void> {
    const barn = await this.quickCreateBarn.open();
    if (barn) {
      this.barns.update((b) => [...b, barn]);
      this.horseForm.controls.barnId.setValue(barn.id);
    }
  }

  goToClient(): void {
    this.router.navigate(['/clients', this.createdClient()!.id]);
  }

  goToNewHorse(): void {
    this.router.navigate(['/horses/new']);
  }

  goToNewAppointment(): void {
    this.router.navigate(['/appointments/new']);
  }
}
