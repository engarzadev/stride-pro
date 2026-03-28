import { Component, Injectable, inject, signal } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { Barn } from '../../../core/models';
import { BarnsService } from '../../../features/barns/barns.service';

@Injectable({ providedIn: 'root' })
export class QuickCreateBarnService {
  readonly visible = signal(false);
  private resolveFn?: (result: Barn | null) => void;

  open(): Promise<Barn | null> {
    this.visible.set(true);
    return new Promise<Barn | null>((resolve) => {
      this.resolveFn = resolve;
    });
  }

  complete(barn: Barn): void {
    this.visible.set(false);
    this.resolveFn?.(barn);
  }

  cancel(): void {
    this.visible.set(false);
    this.resolveFn?.(null);
  }
}

@Component({
  selector: 'app-quick-create-barn',
  standalone: true,
  imports: [ReactiveFormsModule],
  templateUrl: './quick-create-barn.component.html',
  styleUrls: ['./quick-create-modal.scss'],
})
export class QuickCreateBarnComponent {
  readonly service = inject(QuickCreateBarnService);
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
      next: (barn) => {
        this.saving.set(false);
        this.form.reset({ name: '', contactName: '', address: '', phone: '', email: '', notes: '' });
        this.service.complete(barn);
      },
      error: () => this.saving.set(false),
    });
  }

  onCancel(): void {
    this.form.reset({ name: '', contactName: '', address: '', phone: '', email: '', notes: '' });
    this.service.cancel();
  }
}
