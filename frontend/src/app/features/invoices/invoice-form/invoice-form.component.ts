import {
  Component,
  computed,
  inject,
  OnInit,
  signal,
  ViewChild,
} from '@angular/core';
import {
  FormArray,
  FormBuilder,
  ReactiveFormsModule,
  Validators,
} from '@angular/forms';
import {
  MatAutocompleteModule,
  MatAutocompleteSelectedEvent,
} from '@angular/material/autocomplete';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { MatDatepickerModule } from '@angular/material/datepicker';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatSelect, MatSelectModule } from '@angular/material/select';
import { ActivatedRoute, Router } from '@angular/router';
import { Client, ServiceItem } from '../../../core/models';
import { FormPageComponent } from '../../../shared/components/form-page/form-page.component';
import { LoadingSpinnerComponent } from '../../../shared/components/loading-spinner/loading-spinner.component';
import { QuickCreateClientService } from '../../../shared/components/quick-create/quick-create-client.component';
import { ToastService } from '../../../shared/components/toast/toast.service';
import { CurrencyFormatPipe } from '../../../shared/pipes/currency-format.pipe';
import { ClientsService } from '../../clients/clients.service';
import { SettingsService } from '../../settings/settings.service';
import { InvoicesService } from '../invoices.service';

@Component({
  selector: 'app-invoice-form',
  standalone: true,
  imports: [
    ReactiveFormsModule,
    FormPageComponent,
    LoadingSpinnerComponent,
    CurrencyFormatPipe,
    MatFormFieldModule,
    MatInputModule,
    MatSelectModule,
    MatButtonModule,
    MatIconModule,
    MatCardModule,
    MatDatepickerModule,
    MatAutocompleteModule,
  ],
  templateUrl: './invoice-form.component.html',
  styleUrls: ['./invoice-form.component.scss'],
})
export class InvoiceFormComponent implements OnInit {
  private readonly fb = inject(FormBuilder);
  private readonly route = inject(ActivatedRoute);
  private readonly router = inject(Router);
  private readonly invoicesService = inject(InvoicesService);
  private readonly clientsService = inject(ClientsService);
  private readonly settingsService = inject(SettingsService);
  private readonly toast = inject(ToastService);
  private readonly quickCreateClient = inject(QuickCreateClientService);

  @ViewChild('clientSelect') clientSelect!: MatSelect;

  readonly loading = signal(false);
  readonly saving = signal(false);
  readonly isEdit = signal(false);
  readonly clients = signal<Client[]>([]);
  readonly serviceItems = signal<ServiceItem[]>([]);
  readonly descriptionQuery = signal('');
  activeItemIndex = 0;
  private invoiceId = '';

  readonly filteredServiceItems = computed(() => {
    const q = this.descriptionQuery().toLowerCase().trim();
    if (!q) return this.serviceItems();
    return this.serviceItems().filter((s) => s.name.toLowerCase().includes(q));
  });

  readonly form = this.fb.nonNullable.group({
    clientId: ['', [Validators.required]],
    date: [null as Date | null, [Validators.required]],
    dueDate: [null as Date | null, [Validators.required]],
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
    this.settingsService.getServiceItems().subscribe({
      next: (items) => this.serviceItems.set(items),
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
            date: invoice.createdAt ? new Date(invoice.createdAt) : null,
            dueDate: invoice.dueDate ? new Date(invoice.dueDate) : null,
            status: invoice.status,
            notes: invoice.notes,
          });

          this.items.clear();
          if (invoice.items && invoice.items.length > 0) {
            invoice.items.forEach((item) => {
              this.items.push(
                this.fb.nonNullable.group({
                  description: [item.description, [Validators.required]],
                  quantity: [
                    item.quantity,
                    [Validators.required, Validators.min(1)],
                  ],
                  unitPrice: [
                    item.unitPrice,
                    [Validators.required, Validators.min(0)],
                  ],
                  notes: [item.notes ?? ''],
                }),
              );
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
      notes: [''],
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

  onDescriptionInput(value: string): void {
    this.descriptionQuery.set(value);
    this.recalculate();
  }

  onServiceSelected(event: MatAutocompleteSelectedEvent, index: number): void {
    const selected = this.serviceItems().find(
      (s) => s.name === event.option.value,
    );
    if (selected) {
      const item = this.items.at(index);
      item.patchValue({
        description: selected.name,
        unitPrice: selected.defaultPrice,
      });
      this.recalculate();
    }
    this.descriptionQuery.set('');
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
    const toDateTime = (d: Date | null) =>
      d ? `${d.toISOString().substring(0, 10)}T00:00:00Z` : '';
    const data = {
      clientId: formValue.clientId,
      createdAt: toDateTime(formValue.date),
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
        notes: item.notes,
      })),
    };

    const request$ = this.isEdit()
      ? this.invoicesService.update(this.invoiceId, data)
      : this.invoicesService.create(data);

    request$.subscribe({
      next: () => {
        this.toast.success(
          this.isEdit()
            ? 'Invoice updated successfully'
            : 'Invoice created successfully',
        );
        this.router.navigate(['/invoices']);
      },
      error: () => this.saving.set(false),
    });
  }
}
