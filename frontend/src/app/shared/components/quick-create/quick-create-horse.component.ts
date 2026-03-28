import {
  Component,
  Injectable,
  OnInit,
  computed,
  effect,
  inject,
  signal,
} from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import {
  MAT_DIALOG_DATA,
  MatDialog,
  MatDialogModule,
  MatDialogRef,
} from '@angular/material/dialog';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';
import { firstValueFrom } from 'rxjs';
import { Barn, Client, Horse } from '../../../core/models';
import { BarnsService } from '../../../features/barns/barns.service';
import { ClientsService } from '../../../features/clients/clients.service';
import { HorsesService } from '../../../features/horses/horses.service';
import { BreedAutocompleteComponent } from '../breed-autocomplete/breed-autocomplete.component';
import { QuickCreateBarnService } from './quick-create-barn.component';

export interface QuickCreateHorseOptions {
  clientId?: string;
}

@Injectable({ providedIn: 'root' })
export class QuickCreateHorseService {
  private readonly dialog = inject(MatDialog);

  async open(options: QuickCreateHorseOptions = {}): Promise<Horse | null> {
    const ref = this.dialog.open(QuickCreateHorseComponent, {
      width: '650px',
      data: options,
    });
    return (await firstValueFrom(ref.afterClosed())) ?? null;
  }
}

@Component({
  selector: 'app-quick-create-horse',
  standalone: true,
  imports: [
    ReactiveFormsModule,
    MatDialogModule,
    MatFormFieldModule,
    MatInputModule,
    MatSelectModule,
    MatButtonModule,
    MatIconModule,
    BreedAutocompleteComponent,
  ],
  template: `
    <h2 mat-dialog-title>New Horse</h2>
    <mat-dialog-content>
      <form [formGroup]="form" id="qc-horse-form" (ngSubmit)="onSubmit()">
        <mat-form-field appearance="outline">
          <mat-label>Name</mat-label>
          <input matInput formControlName="name" />
          @if (form.controls.name.errors?.['required']) {
            <mat-error>Horse name is required.</mat-error>
          }
        </mat-form-field>
        <div class="form-row">
          <mat-form-field appearance="outline">
            <mat-label>Breed</mat-label>
            <app-breed-autocomplete formControlName="breed" />
          </mat-form-field>
          <mat-form-field appearance="outline">
            <mat-label>Color</mat-label>
            <input matInput formControlName="color" />
          </mat-form-field>
        </div>
        <div class="form-row">
          <mat-form-field appearance="outline">
            <mat-label>Age</mat-label>
            <input matInput type="number" formControlName="age" min="0" />
          </mat-form-field>
          <mat-form-field appearance="outline">
            <mat-label>Gender</mat-label>
            <mat-select formControlName="gender">
              <mat-option value="">Select gender</mat-option>
              <mat-option value="mare">Mare</mat-option>
              <mat-option value="stallion">Stallion</mat-option>
              <mat-option value="gelding">Gelding</mat-option>
            </mat-select>
          </mat-form-field>
        </div>
        <mat-form-field appearance="outline">
          <mat-label>Weight (lbs)</mat-label>
          <input matInput type="number" formControlName="weight" min="0" />
        </mat-form-field>
        <div class="form-row">
          <mat-form-field appearance="outline">
            <mat-label>Client</mat-label>
            @if (preselectedClient(); as client) {
              <input
                matInput
                [value]="client.firstName + ' ' + client.lastName"
                readonly
              />
            } @else {
              <mat-select formControlName="clientId">
                <mat-option [value]="0">Select a client</mat-option>
                @for (c of clients(); track c.id) {
                  <mat-option [value]="c.id"
                    >{{ c.firstName }} {{ c.lastName }}</mat-option
                  >
                }
              </mat-select>
              @if (form.controls.clientId.errors) {
                <mat-error>Please select a client.</mat-error>
              }
            }
          </mat-form-field>
          <mat-form-field appearance="outline">
            <mat-label>Barn</mat-label>
            <mat-select formControlName="barnId">
              <mat-option [value]="null">No barn assigned</mat-option>
              @for (barn of barns(); track barn.id) {
                <mat-option [value]="barn.id">{{ barn.name }}</mat-option>
              }
            </mat-select>
            <button
              matSuffix
              mat-icon-button
              type="button"
              (click)="openCreateBarn()"
              title="Create new barn"
            >
              <mat-icon>add</mat-icon>
            </button>
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
        form="qc-horse-form"
        type="submit"
        [disabled]="saving()"
      >
        @if (saving()) {
          Saving...
        } @else {
          Create Horse
        }
      </button>
    </mat-dialog-actions>
  `,
})
export class QuickCreateHorseComponent implements OnInit {
  readonly dialogRef = inject(MatDialogRef<QuickCreateHorseComponent>);
  readonly options = inject<QuickCreateHorseOptions>(MAT_DIALOG_DATA);
  private readonly horsesService = inject(HorsesService);
  private readonly clientsService = inject(ClientsService);
  private readonly barnsService = inject(BarnsService);
  private readonly quickCreateBarn = inject(QuickCreateBarnService);
  private readonly fb = inject(FormBuilder);

  readonly saving = signal(false);
  readonly clients = signal<Client[]>([]);
  readonly barns = signal<Barn[]>([]);
  readonly preselectedClient = computed(() => {
    const clientId = this.options?.clientId;
    return clientId
      ? (this.clients().find((c) => c.id === clientId) ?? null)
      : null;
  });

  readonly form = this.fb.nonNullable.group({
    name: ['', [Validators.required]],
    breed: [''],
    age: [0],
    gender: [''],
    color: [''],
    weight: [0],
    notes: [''],
    clientId: ['', [Validators.required]],
    barnId: [null as string | null],
  });

  constructor() {
    effect(() => {
      const clientId = this.options?.clientId;
      if (clientId) {
        this.form.controls.clientId.setValue(clientId);
      }
    });
  }

  ngOnInit(): void {
    this.clientsService.getAll().subscribe((c) => this.clients.set(c));
    this.barnsService.getAll().subscribe((b) => this.barns.set(b));
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
    this.horsesService.create(this.form.getRawValue()).subscribe({
      next: (horse) => this.dialogRef.close(horse),
      error: () => this.saving.set(false),
    });
  }
}
