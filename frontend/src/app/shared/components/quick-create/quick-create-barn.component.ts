import { Component, Injectable, inject, signal } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { MatDialog, MatDialogModule, MatDialogRef } from '@angular/material/dialog';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { firstValueFrom } from 'rxjs';
import { Barn } from '../../../core/models';
import { BarnsService } from '../../../features/barns/barns.service';

@Injectable({ providedIn: 'root' })
export class QuickCreateBarnService {
  private readonly dialog = inject(MatDialog);

  async open(): Promise<Barn | null> {
    const ref = this.dialog.open(QuickCreateBarnComponent, { width: '600px' });
    return (await firstValueFrom(ref.afterClosed())) ?? null;
  }
}

@Component({
  selector: 'app-quick-create-barn',
  standalone: true,
  imports: [ReactiveFormsModule, MatDialogModule, MatFormFieldModule, MatInputModule, MatButtonModule],
  template: `
    <h2 mat-dialog-title>New Barn</h2>
    <mat-dialog-content>
      <form [formGroup]="form" id="qc-barn-form" (ngSubmit)="onSubmit()">
        <mat-form-field appearance="outline">
          <mat-label>Name *</mat-label>
          <input matInput formControlName="name">
          @if (form.controls.name.errors?.['required']) {
            <mat-error>Barn name is required.</mat-error>
          }
        </mat-form-field>
        <mat-form-field appearance="outline">
          <mat-label>Contact Name</mat-label>
          <input matInput formControlName="contactName">
        </mat-form-field>
        <mat-form-field appearance="outline">
          <mat-label>Address</mat-label>
          <input matInput formControlName="address">
        </mat-form-field>
        <div class="form-row">
          <mat-form-field appearance="outline">
            <mat-label>Phone</mat-label>
            <input matInput formControlName="phone">
          </mat-form-field>
          <mat-form-field appearance="outline">
            <mat-label>Email</mat-label>
            <input matInput type="email" formControlName="email">
            @if (form.controls.email.errors?.['email']) {
              <mat-error>Please enter a valid email.</mat-error>
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
      <button mat-raised-button color="primary" form="qc-barn-form" type="submit" [disabled]="saving()">
        @if (saving()) { Saving... } @else { Create Barn }
      </button>
    </mat-dialog-actions>
  `,
})
export class QuickCreateBarnComponent {
  readonly dialogRef = inject(MatDialogRef<QuickCreateBarnComponent>);
  private readonly barnsService = inject(BarnsService);
  private readonly fb = inject(FormBuilder);
  readonly saving = signal(false);

  readonly form = this.fb.nonNullable.group({
    name: ['', [Validators.required]],
    contactName: [''],
    address: [''],
    phone: [''],
    email: ['', [Validators.email]],
    notes: [''],
  });

  onSubmit(): void {
    if (this.form.invalid) {
      this.form.markAllAsTouched();
      return;
    }
    this.saving.set(true);
    this.barnsService.create(this.form.getRawValue()).subscribe({
      next: (barn) => this.dialogRef.close(barn),
      error: () => this.saving.set(false),
    });
  }
}
