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
import { Observable, forkJoin, firstValueFrom, of } from 'rxjs';
import { Barn, Horse } from '../../../core/models';
import { SubscriptionService } from '../../../core/services/subscription.service';
import { BarnsService } from '../../../features/barns/barns.service';
import { HorsesService } from '../../../features/horses/horses.service';
import { ToastService } from '../toast/toast.service';
import { HorseMultiselectAutocompleteComponent } from '../horse-multiselect-autocomplete/horse-multiselect-autocomplete.component';

@Injectable({ providedIn: 'root' })
export class QuickCreateBarnService {
  private readonly dialog = inject(MatDialog);
  private readonly subscriptionService = inject(SubscriptionService);
  private readonly toast = inject(ToastService);

  async open(): Promise<Barn | null> {
    await firstValueFrom(this.subscriptionService.load());
    if (!this.subscriptionService.hasFeature('barn_management')) {
      this.toast.error('Barn management requires a paid plan');
      return null;
    }
    const ref = this.dialog.open(QuickCreateBarnComponent, { width: '600px' });
    return (await firstValueFrom(ref.afterClosed())) ?? null;
  }
}

@Component({
  selector: 'app-quick-create-barn',
  standalone: true,
  imports: [
    ReactiveFormsModule,
    MatDialogModule,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule,
    HorseMultiselectAutocompleteComponent,
  ],
  template: `
    <h2 mat-dialog-title>New Barn</h2>
    <mat-dialog-content>
      <form [formGroup]="form" id="qc-barn-form" (ngSubmit)="onSubmit()">
        <mat-form-field appearance="outline">
          <mat-label>Name</mat-label>
          <input matInput formControlName="name" />
          @if (form.controls.name.errors?.['required']) {
            <mat-error>Barn name is required.</mat-error>
          }
        </mat-form-field>
        <mat-form-field appearance="outline">
          <mat-label>Contact Name</mat-label>
          <input matInput formControlName="contactName" />
        </mat-form-field>
        <mat-form-field appearance="outline">
          <mat-label>Address</mat-label>
          <input matInput formControlName="address" />
        </mat-form-field>
        <div class="form-row">
          <mat-form-field appearance="outline">
            <mat-label>Phone</mat-label>
            <input matInput formControlName="phone" />
          </mat-form-field>
          <mat-form-field appearance="outline">
            <mat-label>Email</mat-label>
            <input matInput type="email" formControlName="email" />
            @if (form.controls.email.errors?.['email']) {
              <mat-error>Please enter a valid email.</mat-error>
            }
          </mat-form-field>
        </div>
        <mat-form-field appearance="outline">
          <mat-label>Notes</mat-label>
          <textarea matInput formControlName="notes" rows="2"></textarea>
        </mat-form-field>
        <app-horse-multiselect-autocomplete
          [horses]="allHorses()"
          [selectedHorses]="selectedHorses()"
          (selectedHorsesChange)="selectedHorses.set($event)"
        />
      </form>
    </mat-dialog-content>
    <mat-dialog-actions align="end">
      <button mat-stroked-button (click)="dialogRef.close(null)">Cancel</button>
      <button
        mat-raised-button
        color="primary"
        form="qc-barn-form"
        type="submit"
        [disabled]="saving()"
      >
        @if (saving()) {
          Saving...
        } @else {
          Create Barn
        }
      </button>
    </mat-dialog-actions>
  `,
})
export class QuickCreateBarnComponent implements OnInit {
  readonly dialogRef = inject(MatDialogRef<QuickCreateBarnComponent>);
  private readonly barnsService = inject(BarnsService);
  private readonly horsesService = inject(HorsesService);
  private readonly fb = inject(FormBuilder);
  readonly saving = signal(false);
  readonly allHorses = signal<Horse[]>([]);
  readonly selectedHorses = signal<Horse[]>([]);

  readonly form = this.fb.nonNullable.group({
    name: ['', [Validators.required]],
    contactName: [''],
    address: [''],
    phone: [''],
    email: ['', [Validators.email]],
    notes: [''],
  });

  ngOnInit(): void {
    this.horsesService.getAll().subscribe((horses) => this.allHorses.set(horses));
  }

  onSubmit(): void {
    if (this.form.invalid) {
      this.form.markAllAsTouched();
      return;
    }
    this.saving.set(true);
    this.barnsService.create(this.form.getRawValue()).subscribe({
      next: (barn) => {
        const horses = this.selectedHorses();
        const updates = horses.map((horse) => this.horsesService.update(horse.id, { ...horse, barnId: barn.id }));
        const assign$: Observable<unknown> = updates.length > 0 ? forkJoin(updates) : of(null);
        assign$.subscribe({
          next: () => this.dialogRef.close(barn),
          error: () => this.dialogRef.close(barn),
        });
      },
      error: () => this.saving.set(false),
    });
  }
}
