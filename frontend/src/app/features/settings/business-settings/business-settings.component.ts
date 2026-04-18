import { Component, inject, OnInit, signal } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { MatCardModule } from '@angular/material/card';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { SettingsService } from '../settings.service';
import { ToastService } from '../../../shared/components/toast/toast.service';
import { LoadingSpinnerComponent } from '../../../shared/components/loading-spinner/loading-spinner.component';

@Component({
  selector: 'app-business-settings',
  standalone: true,
  imports: [ReactiveFormsModule, LoadingSpinnerComponent, MatCardModule, MatFormFieldModule, MatInputModule, MatButtonModule],
  templateUrl: './business-settings.component.html',
  styleUrls: ['./business-settings.component.scss'],
})
export class BusinessSettingsComponent implements OnInit {
  private readonly fb = inject(FormBuilder);
  private readonly settingsService = inject(SettingsService);
  private readonly toast = inject(ToastService);

  readonly loading = signal(true);
  readonly saving = signal(false);

  readonly form = this.fb.nonNullable.group({
    businessName: [''],
    email: ['', [Validators.email]],
    phone: [''],
    address: [''],
    invoiceMessage: [''],
  });

  ngOnInit(): void {
    this.settingsService.getBusinessSettings().subscribe({
      next: (bs) => {
        this.form.patchValue({
          businessName: bs.businessName,
          email: bs.email,
          phone: bs.phone,
          address: bs.address,
          invoiceMessage: bs.invoiceMessage,
        });
        this.loading.set(false);
      },
      error: () => this.loading.set(false),
    });
  }

  onSave(): void {
    if (this.form.invalid) {
      this.form.markAllAsTouched();
      return;
    }
    this.saving.set(true);
    const value = this.form.getRawValue();
    this.settingsService.saveBusinessSettings({
      businessName: value.businessName,
      email: value.email,
      phone: value.phone,
      address: value.address,
      invoiceMessage: value.invoiceMessage,
    }).subscribe({
      next: () => {
        this.toast.success('Business settings saved');
        this.saving.set(false);
      },
      error: () => this.saving.set(false),
    });
  }
}
