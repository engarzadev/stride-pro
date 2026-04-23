import { Component, computed, inject, signal } from '@angular/core';
import { toSignal } from '@angular/core/rxjs-interop';
import {
  AbstractControl,
  FormControl,
  FormGroup,
  ReactiveFormsModule,
  ValidationErrors,
  Validators,
} from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { AuthService } from '../../../core/services/auth.service';
import { ToastService } from '../../../shared/components/toast/toast.service';

function passwordsMatch(group: AbstractControl): ValidationErrors | null {
  const newPw = group.get('newPassword')?.value;
  const confirm = group.get('confirmPassword')?.value;
  return newPw && confirm && newPw !== confirm ? { mismatch: true } : null;
}

@Component({
  selector: 'app-general-settings',
  standalone: true,
  imports: [
    ReactiveFormsModule,
    MatCardModule,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule,
    MatIconModule,
  ],
  templateUrl: './general-settings.component.html',
  styleUrls: ['./general-settings.component.scss'],
})
export class GeneralSettingsComponent {
  private readonly authService = inject(AuthService);
  private readonly toast = inject(ToastService);
  private readonly currentUser = toSignal(this.authService.currentUser$);

  readonly isOwner = computed(() => this.currentUser()?.role === 'owner');
  readonly isProfessional = computed(() => this.currentUser()?.role !== 'owner');
  readonly subscriptionTier = computed(() => this.currentUser()?.subscriptionTier ?? 'free');

  readonly savingProfile = signal(false);
  readonly savingPassword = signal(false);
  readonly showCurrentPassword = signal(false);
  readonly showNewPassword = signal(false);

  readonly profileForm = new FormGroup({
    firstName: new FormControl('', { nonNullable: true, validators: [Validators.required] }),
    lastName: new FormControl('', { nonNullable: true, validators: [Validators.required] }),
    email: new FormControl('', {
      nonNullable: true,
      validators: [Validators.required, Validators.email],
    }),
  });

  readonly passwordForm = new FormGroup(
    {
      currentPassword: new FormControl('', {
        nonNullable: true,
        validators: [Validators.required],
      }),
      newPassword: new FormControl('', {
        nonNullable: true,
        validators: [Validators.required, Validators.minLength(8)],
      }),
      confirmPassword: new FormControl('', {
        nonNullable: true,
        validators: [Validators.required],
      }),
    },
    { validators: passwordsMatch },
  );

  readonly tierConfig = computed(() => {
    const configs: Record<string, { label: string; description: string; next?: string }> = {
      free: {
        label: 'Free',
        description: 'Up to 10 clients and 20 horses with basic scheduling.',
        next: 'base',
      },
      base: {
        label: 'Base',
        description: 'Unlimited clients and horses with full scheduling.',
        next: 'trainer',
      },
      trainer: {
        label: 'Trainer',
        description: 'Multi-horse sessions and advanced reporting.',
        next: 'enterprise',
      },
      enterprise: {
        label: 'Enterprise',
        description: 'SMS notifications, API access, and custom branding.',
      },
    };
    return configs[this.subscriptionTier()] ?? configs['free'];
  });

  constructor() {
    const user = this.currentUser();
    if (user) {
      this.profileForm.patchValue({
        firstName: user.firstName,
        lastName: user.lastName,
        email: user.email,
      });
    }
  }

  saveProfile(): void {
    if (this.profileForm.invalid) return;
    this.savingProfile.set(true);
    const { firstName, lastName, email } = this.profileForm.getRawValue();
    this.authService.updateProfile(firstName, lastName, email).subscribe({
      next: () => {
        this.toast.success('Profile updated.');
        this.savingProfile.set(false);
      },
      error: () => this.savingProfile.set(false),
    });
  }

  updatePassword(): void {
    if (this.passwordForm.invalid) return;
    this.savingPassword.set(true);
    const { currentPassword, newPassword } = this.passwordForm.getRawValue();
    this.authService.changePassword(currentPassword, newPassword).subscribe({
      next: () => {
        this.toast.success('Password updated.');
        this.passwordForm.reset();
        this.savingPassword.set(false);
      },
      error: () => this.savingPassword.set(false),
    });
  }
}
