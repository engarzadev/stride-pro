import { Component, Injectable, inject, signal } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { MatDialog, MatDialogModule, MatDialogRef } from '@angular/material/dialog';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { firstValueFrom } from 'rxjs';
import { Client } from '../../../core/models';
import { ClientsService } from '../../../features/clients/clients.service';

@Injectable({ providedIn: 'root' })
export class QuickCreateClientService {
  private readonly dialog = inject(MatDialog);

  async open(): Promise<Client | null> {
    const ref = this.dialog.open(QuickCreateClientComponent, { width: '600px' });
    return (await firstValueFrom(ref.afterClosed())) ?? null;
  }
}

@Component({
  selector: 'app-quick-create-client',
  standalone: true,
  imports: [ReactiveFormsModule, MatDialogModule, MatFormFieldModule, MatInputModule, MatButtonModule],
  template: `
    <h2 mat-dialog-title>New Client</h2>
    <mat-dialog-content>
      <form [formGroup]="form" id="qc-client-form" (ngSubmit)="onSubmit()">
        <div class="form-row">
          <mat-form-field appearance="outline">
            <mat-label>First Name *</mat-label>
            <input matInput formControlName="firstName">
            @if (form.controls.firstName.errors?.['required']) {
              <mat-error>First name is required.</mat-error>
            }
          </mat-form-field>
          <mat-form-field appearance="outline">
            <mat-label>Last Name *</mat-label>
            <input matInput formControlName="lastName">
            @if (form.controls.lastName.errors?.['required']) {
              <mat-error>Last name is required.</mat-error>
            }
          </mat-form-field>
        </div>
        <div class="form-row">
          <mat-form-field appearance="outline">
            <mat-label>Email</mat-label>
            <input matInput type="email" formControlName="email">
            @if (form.controls.email.errors?.['email']) {
              <mat-error>Please enter a valid email.</mat-error>
            }
          </mat-form-field>
          <mat-form-field appearance="outline">
            <mat-label>Phone</mat-label>
            <input matInput formControlName="phone">
          </mat-form-field>
        </div>
        <mat-form-field appearance="outline">
          <mat-label>Address</mat-label>
          <input matInput formControlName="address">
        </mat-form-field>
        <mat-form-field appearance="outline">
          <mat-label>Notes</mat-label>
          <textarea matInput formControlName="notes" rows="2"></textarea>
        </mat-form-field>
      </form>
    </mat-dialog-content>
    <mat-dialog-actions align="end">
      <button mat-stroked-button (click)="dialogRef.close(null)">Cancel</button>
      <button mat-raised-button color="primary" form="qc-client-form" type="submit" [disabled]="saving()">
        @if (saving()) { Saving... } @else { Create Client }
      </button>
    </mat-dialog-actions>
  `,
})
export class QuickCreateClientComponent {
  readonly dialogRef = inject(MatDialogRef<QuickCreateClientComponent>);
  private readonly clientsService = inject(ClientsService);
  private readonly fb = inject(FormBuilder);
  readonly saving = signal(false);

  readonly form = this.fb.nonNullable.group({
    firstName: ['', [Validators.required]],
    lastName: ['', [Validators.required]],
    email: ['', [Validators.email]],
    phone: [''],
    address: [''],
    notes: [''],
  });

  onSubmit(): void {
    if (this.form.invalid) {
      this.form.markAllAsTouched();
      return;
    }
    this.saving.set(true);
    this.clientsService.create(this.form.getRawValue()).subscribe({
      next: (client) => this.dialogRef.close(client),
      error: () => this.saving.set(false),
    });
  }
}
