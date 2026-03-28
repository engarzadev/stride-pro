import { Component, Injectable, inject, signal } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { Client } from '../../../core/models';
import { ClientsService } from '../../../features/clients/clients.service';

@Injectable({ providedIn: 'root' })
export class QuickCreateClientService {
  readonly visible = signal(false);
  private resolveFn?: (result: Client | null) => void;

  open(): Promise<Client | null> {
    this.visible.set(true);
    return new Promise<Client | null>((resolve) => {
      this.resolveFn = resolve;
    });
  }

  complete(client: Client): void {
    this.visible.set(false);
    this.resolveFn?.(client);
  }

  cancel(): void {
    this.visible.set(false);
    this.resolveFn?.(null);
  }
}

@Component({
  selector: 'app-quick-create-client',
  standalone: true,
  imports: [ReactiveFormsModule],
  templateUrl: './quick-create-client.component.html',
  styleUrls: ['./quick-create-modal.scss'],
})
export class QuickCreateClientComponent {
  readonly service = inject(QuickCreateClientService);
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
      next: (client) => {
        this.saving.set(false);
        this.form.reset({ firstName: '', lastName: '', email: '', phone: '', address: '', notes: '' });
        this.service.complete(client);
      },
      error: () => this.saving.set(false),
    });
  }

  onCancel(): void {
    this.form.reset({ firstName: '', lastName: '', email: '', phone: '', address: '', notes: '' });
    this.service.cancel();
  }
}
