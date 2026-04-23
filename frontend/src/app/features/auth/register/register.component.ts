import { Component, computed, signal } from '@angular/core';
import { inject } from '@angular/core';
import { Router, RouterLink } from '@angular/router';
import { FormControl, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { AuthService } from '../../../core/services/auth.service';
import { ToastService } from '../../../shared/components/toast/toast.service';

@Component({
  selector: 'app-register',
  standalone: true,
  imports: [ReactiveFormsModule, RouterLink, MatFormFieldModule, MatInputModule, MatButtonModule, MatIconModule],
  templateUrl: './register.component.html',
  styleUrls: ['./register.component.scss'],
})
export class RegisterComponent {
  private readonly authService = inject(AuthService);
  private readonly router = inject(Router);
  private readonly toast = inject(ToastService);

  readonly loading = signal(false);
  readonly showPassword = signal(false);
  readonly selectedAccountType = signal<'owner' | 'professional' | null>(null);
  readonly showAccountTypeError = signal(false);

  // Per-field updateOn:'blur' so errors appear as user leaves each field
  readonly form = new FormGroup({
    firstName: new FormControl('', { nonNullable: true, validators: [Validators.required], updateOn: 'blur' }),
    lastName: new FormControl('', { nonNullable: true, validators: [Validators.required], updateOn: 'blur' }),
    email: new FormControl('', { nonNullable: true, validators: [Validators.required, Validators.email], updateOn: 'blur' }),
    password: new FormControl('', { nonNullable: true, validators: [Validators.required, Validators.minLength(8)] }),
  });

  readonly passwordValue = signal('');
  readonly passwordRequirements = computed(() => {
    const pw = this.passwordValue();
    return [
      { label: 'At least 8 characters', met: pw.length >= 8 },
      { label: 'One uppercase letter', met: /[A-Z]/.test(pw) },
      { label: 'One lowercase letter', met: /[a-z]/.test(pw) },
      { label: 'One number', met: /\d/.test(pw) },
    ];
  });
  readonly allRequirementsMet = computed(() => this.passwordRequirements().every(r => r.met));
  readonly passwordTouched = signal(false);

  selectAccountType(type: 'owner' | 'professional'): void {
    this.selectedAccountType.set(type);
    this.showAccountTypeError.set(false);
  }

  onPasswordInput(event: Event): void {
    this.passwordValue.set((event.target as HTMLInputElement).value);
  }

  onSubmit(): void {
    if (!this.selectedAccountType()) {
      this.showAccountTypeError.set(true);
      return;
    }

    this.form.markAllAsTouched();
    if (this.form.invalid) {
      return;
    }

    this.loading.set(true);
    const { email, password, firstName, lastName } = this.form.getRawValue();
    this.authService.register(email, password, firstName, lastName, this.selectedAccountType()!).subscribe({
      next: () => {
        this.toast.success('Account created successfully!');
        this.router.navigate(['/dashboard']);
      },
      error: () => {
        this.loading.set(false);
      },
    });
  }
}
