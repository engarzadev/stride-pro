import { Component, inject, OnInit, signal } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { MatCardModule } from '@angular/material/card';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatTableModule } from '@angular/material/table';
import { CurrencyFormatPipe } from '../../../shared/pipes/currency-format.pipe';
import { SettingsService } from '../settings.service';
import { ToastService } from '../../../shared/components/toast/toast.service';
import { ConfirmDialogService } from '../../../shared/components/confirm-dialog/confirm-dialog.component';
import { LoadingSpinnerComponent } from '../../../shared/components/loading-spinner/loading-spinner.component';
import { ServiceItem } from '../../../core/models';

@Component({
  selector: 'app-service-catalog',
  standalone: true,
  imports: [ReactiveFormsModule, LoadingSpinnerComponent, CurrencyFormatPipe, MatCardModule, MatFormFieldModule, MatInputModule, MatButtonModule, MatIconModule, MatTableModule],
  templateUrl: './service-catalog.component.html',
  styleUrls: ['./service-catalog.component.scss'],
})
export class ServiceCatalogComponent implements OnInit {
  private readonly fb = inject(FormBuilder);
  private readonly settingsService = inject(SettingsService);
  private readonly toast = inject(ToastService);
  private readonly confirmDialog = inject(ConfirmDialogService);

  readonly loading = signal(true);
  readonly saving = signal(false);
  readonly items = signal<ServiceItem[]>([]);
  readonly editingId = signal<string | null>(null);

  readonly displayedColumns = ['name', 'defaultPrice', 'actions'];

  readonly addForm = this.fb.nonNullable.group({
    name: ['', [Validators.required]],
    defaultPrice: [0, [Validators.required, Validators.min(0)]],
  });

  readonly editForm = this.fb.nonNullable.group({
    name: ['', [Validators.required]],
    defaultPrice: [0, [Validators.required, Validators.min(0)]],
  });

  ngOnInit(): void {
    this.settingsService.getServiceItems().subscribe({
      next: (items) => {
        this.items.set(items);
        this.loading.set(false);
      },
      error: () => this.loading.set(false),
    });
  }

  onAdd(): void {
    if (this.addForm.invalid) {
      this.addForm.markAllAsTouched();
      return;
    }
    this.saving.set(true);
    const value = this.addForm.getRawValue();
    this.settingsService.createServiceItem({ name: value.name, defaultPrice: value.defaultPrice }).subscribe({
      next: (item) => {
        this.items.update((prev) => [...prev, item]);
        this.addForm.reset({ name: '', defaultPrice: 0 });
        this.saving.set(false);
        this.toast.success('Service added');
      },
      error: () => this.saving.set(false),
    });
  }

  startEdit(item: ServiceItem): void {
    this.editingId.set(item.id);
    this.editForm.setValue({ name: item.name, defaultPrice: item.defaultPrice });
  }

  cancelEdit(): void {
    this.editingId.set(null);
  }

  onSaveEdit(item: ServiceItem): void {
    if (this.editForm.invalid) {
      this.editForm.markAllAsTouched();
      return;
    }
    const value = this.editForm.getRawValue();
    this.settingsService.updateServiceItem(item.id, { name: value.name, defaultPrice: value.defaultPrice }).subscribe({
      next: (updated) => {
        this.items.update((prev) => prev.map((i) => (i.id === updated.id ? updated : i)));
        this.editingId.set(null);
        this.toast.success('Service updated');
      },
    });
  }

  async onDelete(item: ServiceItem): Promise<void> {
    const confirmed = await this.confirmDialog.confirm({
      title: 'Delete Service',
      message: `Remove "${item.name}" from your catalog?`,
      confirmText: 'Delete',
      confirmClass: 'btn-danger',
    });
    if (!confirmed) return;

    this.settingsService.deleteServiceItem(item.id).subscribe({
      next: () => {
        this.items.update((prev) => prev.filter((i) => i.id !== item.id));
        this.toast.success('Service removed');
      },
    });
  }
}
