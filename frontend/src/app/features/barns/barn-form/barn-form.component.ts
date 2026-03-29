import { Component, inject, OnInit, signal } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { FormPageComponent } from '../../../shared/components/form-page/form-page.component';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatCardModule } from '@angular/material/card';
import { SubscriptionService } from '../../../core/services/subscription.service';
import { BarnsService } from '../barns.service';
import { ToastService } from '../../../shared/components/toast/toast.service';
import { LoadingSpinnerComponent } from '../../../shared/components/loading-spinner/loading-spinner.component';

@Component({
  selector: 'app-barn-form',
  standalone: true,
  imports: [ReactiveFormsModule, FormPageComponent, LoadingSpinnerComponent, MatFormFieldModule, MatInputModule, MatButtonModule, MatIconModule, MatCardModule],
  templateUrl: './barn-form.component.html',
  styleUrls: ['./barn-form.component.scss'],
})
export class BarnFormComponent implements OnInit {
  private readonly fb = inject(FormBuilder);
  private readonly route = inject(ActivatedRoute);
  private readonly router = inject(Router);
  private readonly barnsService = inject(BarnsService);
  private readonly toast = inject(ToastService);
  private readonly subscriptionService = inject(SubscriptionService);

  readonly loading = signal(false);
  readonly saving = signal(false);
  readonly isEdit = signal(false);
  private barnId = '';

  readonly form = this.fb.nonNullable.group({
    name: ['', [Validators.required]],
    contactName: [''],
    address: [''],
    phone: [''],
    email: ['', [Validators.email]],
    notes: [''],
  });

  ngOnInit(): void {
    const id = this.route.snapshot.paramMap.get('id');
    if (!id) {
      // Creating a new barn — check feature access
      this.subscriptionService.load().subscribe(() => {
        if (!this.subscriptionService.hasFeature('barn_management')) {
          this.toast.error('Barn management requires a paid plan');
          this.router.navigate(['/barns']);
        }
      });
    }
    if (id) {
      this.isEdit.set(true);
      this.barnId = id;
      this.loading.set(true);
      this.barnsService.getById(this.barnId).subscribe({
        next: (barn) => {
          this.form.patchValue({
            name: barn.name,
            contactName: barn.contactName,
            address: barn.address,
            phone: barn.phone,
            email: barn.email,
            notes: barn.notes,
          });
          this.loading.set(false);
        },
        error: () => {
          this.loading.set(false);
          this.router.navigate(['/barns']);
        },
      });
    }
  }

  onSubmit(): void {
    if (this.form.invalid) {
      this.form.markAllAsTouched();
      return;
    }

    this.saving.set(true);
    const data = this.form.getRawValue();

    const request$ = this.isEdit()
      ? this.barnsService.update(this.barnId, data)
      : this.barnsService.create(data);

    request$.subscribe({
      next: () => {
        this.toast.success(this.isEdit() ? 'Barn updated successfully' : 'Barn created successfully');
        this.router.navigate(['/barns']);
      },
      error: () => this.saving.set(false),
    });
  }
}
