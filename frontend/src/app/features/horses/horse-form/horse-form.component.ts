import { Component, computed, inject, OnInit, signal } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { FormPageComponent } from '../../../shared/components/form-page/form-page.component';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';
import { UpgradeFieldPromptComponent } from '../../../shared/components/upgrade-field-prompt/upgrade-field-prompt.component';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatCardModule } from '@angular/material/card';
import { HorsesService } from '../horses.service';
import { ClientsService } from '../../clients/clients.service';
import { BarnsService } from '../../barns/barns.service';
import { Client, Barn, Horse } from '../../../core/models';
import { AuthService } from '../../../core/services/auth.service';
import { SubscriptionService } from '../../../core/services/subscription.service';
import { ToastService } from '../../../shared/components/toast/toast.service';
import { LoadingSpinnerComponent } from '../../../shared/components/loading-spinner/loading-spinner.component';
import { BreedAutocompleteComponent } from '../../../shared/components/breed-autocomplete/breed-autocomplete.component';
import { QuickCreateClientService } from '../../../shared/components/quick-create/quick-create-client.component';
import { QuickCreateBarnService } from '../../../shared/components/quick-create/quick-create-barn.component';

@Component({
  selector: 'app-horse-form',
  standalone: true,
  imports: [ReactiveFormsModule, FormPageComponent, LoadingSpinnerComponent, BreedAutocompleteComponent, MatFormFieldModule, MatInputModule, MatSelectModule, MatButtonModule, MatIconModule, MatCardModule, UpgradeFieldPromptComponent],
  templateUrl: './horse-form.component.html',
  styleUrls: ['./horse-form.component.scss'],
})
export class HorseFormComponent implements OnInit {
  private readonly fb = inject(FormBuilder);
  private readonly route = inject(ActivatedRoute);
  private readonly router = inject(Router);
  private readonly horsesService = inject(HorsesService);
  private readonly clientsService = inject(ClientsService);
  private readonly barnsService = inject(BarnsService);
  private readonly toast = inject(ToastService);
  private readonly quickCreateClient = inject(QuickCreateClientService);
  private readonly quickCreateBarn = inject(QuickCreateBarnService);
  private readonly subscriptionService = inject(SubscriptionService);
  private readonly authService = inject(AuthService);

  readonly loading = signal(false);
  readonly isOwner: boolean;
  readonly saving = signal(false);
  readonly isEdit = signal(false);
  readonly clients = signal<Client[]>([]);
  readonly barns = signal<Barn[]>([]);
  readonly canManageBarns = computed(() => this.subscriptionService.hasFeature('barn_management'));
  private horseId = '';

  readonly form = this.fb.nonNullable.group({
    name: ['', [Validators.required]],
    breed: [''],
    age: [null as number | null],
    gender: [''],
    color: [''],
    weight: [null as number | null],
    notes: [''],
    vetName: [''],
    vetPhone: [''],
    farrierName: [''],
    farrierPhone: [''],
    clientId: ['', [Validators.required]],
    barnId: [null as string | null],
  });

  constructor() {
    const user = this.authService.getStoredUser();
    this.isOwner = user?.role === 'owner';
    if (this.isOwner) {
      this.form.controls.clientId.clearValidators();
      this.form.controls.clientId.updateValueAndValidity();
    }
  }

  ngOnInit(): void {
    if (!this.subscriptionService.hasFeature('barn_management')) {
      this.form.controls.barnId.disable();
    }
    this.clientsService.getAll().subscribe((c) => this.clients.set(c));
    this.barnsService.getAll().subscribe((b) => this.barns.set(b));

    const id = this.route.snapshot.paramMap.get('id');
    if (id) {
      this.isEdit.set(true);
      this.horseId = id;
      this.loading.set(true);
      this.horsesService.getById(this.horseId).subscribe({
        next: (horse) => {
          this.form.patchValue({
            name: horse.name,
            breed: horse.breed,
            age: horse.age,
            gender: horse.gender,
            color: horse.color,
            weight: horse.weight,
            notes: horse.notes,
            vetName: horse.vetName,
            vetPhone: horse.vetPhone,
            farrierName: horse.farrierName,
            farrierPhone: horse.farrierPhone,
            clientId: horse.clientId ?? '',
            barnId: horse.barnId,
          });
          this.loading.set(false);
        },
        error: () => {
          this.loading.set(false);
          this.router.navigate(['/horses']);
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
    const raw = this.form.getRawValue();
    const data = {
      ...raw,
      age: raw.age ?? undefined,
      weight: raw.weight ?? undefined,
      clientId: raw.clientId || null,
      barnId: raw.barnId || null,
    };

    const request$ = this.isEdit()
      ? this.horsesService.update(this.horseId, data as Partial<Horse>)
      : this.horsesService.create(data as Partial<Horse>);

    request$.subscribe({
      next: () => {
        this.toast.success(this.isEdit() ? 'Horse updated successfully' : 'Horse created successfully');
        this.router.navigate(['/horses']);
      },
      error: () => this.saving.set(false),
    });
  }
}
