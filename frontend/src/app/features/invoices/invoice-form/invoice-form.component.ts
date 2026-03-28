import { Component, inject, OnInit, signal, computed } from '@angular/core';
import { ActivatedRoute, Router, RouterLink } from '@angular/router';
import { FormArray, FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { InvoicesService } from '../invoices.service';
import { ClientsService } from '../../clients/clients.service';
import { Client } from '../../../core/models';
import { ToastService } from '../../../shared/components/toast/toast.service';
import { LoadingSpinnerComponent } from '../../../shared/components/loading-spinner/loading-spinner.component';
import { CurrencyFormatPipe } from '../../../shared/pipes/currency-format.pipe';
import { QuickCreateClientService } from '../../../shared/components/quick-create/quick-create-client.component';

@Component({
  selector: 'app-invoice-form',
  standalone: true,
  imports: [ReactiveFormsModule, RouterLink, LoadingSpinnerComponent, CurrencyFormatPipe],
  templateUrl: './invoice-form.component.html',
  styleUrls: ['./invoice-form.component.scss'],
})
export class InvoiceFormComponent implements OnInit {
  private readonly fb = inject(FormBuilder);
  private readonly route = inject(ActivatedRoute);
  private readonly router = inject(Router);
  private readonly invoicesService = inject(InvoicesService);
  private readonly clientsService = inject(ClientsService);
  private readonly toast = inject(ToastService);
  private readonly quickCreateClient = inject(QuickCreateClientService);

  readonly loading = signal(false);
  readonly saving = signal(false);
  readonly isEdit = signal(false);
  readonly clients = signal<Client[]>([]);
  private invoiceId = '';

  readonly form = this.fb.nonNullable.group({
    clientId: ['', [Validators.required]],
    date: ['', [Validators.required]],
    dueDate: ['', [Validators.required]],
    status: ['draft'],
    notes: [''],
    items: this.fb.array([this.createItemGroup()]),
  });

  get items(): FormArray {
    return this.form.get('items') as FormArray;
  }

  readonly subtotal = signal(0);
  readonly tax = signal(0);
  readonly total = signal(0);

  ngOnInit(): void {
    this.clientsService.getAll().subscribe({
      next: (clients) => this.clients.set(clients),
    });

    const id = this.route.snapshot.paramMap.get('id');
    if (id) {
      this.isEdit.set(true);
      this.invoiceId = id;
      this.loading.set(true);
      this.invoicesService.getById(this.invoiceId).subscribe({
        next: (invoice) => {
          this.form.patchValue({
            clientId: invoice.clientId,
            date: invoice.date ? invoice.date.substring(0, 10) : '',
            dueDate: invoice.dueDate ? invoice.dueDate.substring(0, 10) : '',
            status: invoice.status,
            notes: invoice.notes,
          });

          this.items.clear();
          if (invoice.items && invoice.items.length > 0) {
            invoice.items.forEach((item) => {
              this.items.push(this.fb.nonNullable.group({
                description: [item.description, [Validators.required]],
                quantity: [item.quantity, [Validators.required, Validators.min(1)]],
                unitPrice: [item.unitPrice, [Validators.required, Validators.min(0)]],
              }));
            });
          } else {
            this.items.push(this.createItemGroup());
          }

          this.recalculate();
          this.loading.set(false);
        },
        error: () => {
          this.loading.set(false);
          this.router.navigate(['/invoices']);
        },
      });
    }
  }

  async openCreateClient(): Promise<void> {
    const client = await this.quickCreateClient.open();
    if (client) {
      this.clients.update((c) => [...c, client]);
      this.form.controls.clientId.setValue(client.id);
    }
  }

  createItemGroup() {
    return this.fb.nonNullable.group({
      description: ['', [Validators.required]],
      quantity: [1, [Validators.required, Validators.min(1)]],
      unitPrice: [0, [Validators.required, Validators.min(0)]],
    });
  }

  addItem(): void {
    this.items.push(this.createItemGroup());
  }

  removeItem(index: number): void {
    if (this.items.length > 1) {
      this.items.removeAt(index);
      this.recalculate();
    }
  }

  recalculate(): void {
    let subtotal = 0;
    for (let i = 0; i < this.items.length; i++) {
      const item = this.items.at(i);
      const qty = item.get('quantity')?.value || 0;
      const price = item.get('unitPrice')?.value || 0;
      subtotal += qty * price;
    }
    this.subtotal.set(subtotal);
    this.tax.set(0);
    this.total.set(subtotal);
  }

  onSubmit(): void {
    if (this.form.invalid) {
      this.form.markAllAsTouched();
      return;
    }

    this.saving.set(true);
    this.recalculate();

    const formValue = this.form.getRawValue();
    const toDateTime = (d: string) => (d && d.length === 10 ? `${d}T00:00:00Z` : d);
    const data = {
      clientId: formValue.clientId,
      date: toDateTime(formValue.date),
      dueDate: toDateTime(formValue.dueDate),
      status: formValue.status,
      notes: formValue.notes,
      subtotal: this.subtotal(),
      tax: this.tax(),
      total: this.total(),
      items: formValue.items.map((item) => ({
        description: item.description,
        quantity: item.quantity,
        unitPrice: item.unitPrice,
        amount: item.quantity * item.unitPrice,
      })),
    };

    const request$ = this.isEdit()
      ? this.invoicesService.update(this.invoiceId, data)
      : this.invoicesService.create(data);

    request$.subscribe({
      next: () => {
        this.toast.success(this.isEdit() ? 'Invoice updated successfully' : 'Invoice created successfully');
        this.router.navigate(['/invoices']);
      },
      error: () => this.saving.set(false),
    });
  }
}
