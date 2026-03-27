import { Component, inject, OnInit, signal } from '@angular/core';
import { Router } from '@angular/router';
import { InvoicesService } from '../invoices.service';
import { Invoice } from '../../../core/models';
import { PageHeaderComponent } from '../../../shared/components/page-header/page-header.component';
import { DataTableComponent, TableColumn, TableAction } from '../../../shared/components/data-table/data-table.component';
import { LoadingSpinnerComponent } from '../../../shared/components/loading-spinner/loading-spinner.component';
import { ConfirmDialogService } from '../../../shared/components/confirm-dialog/confirm-dialog.component';
import { ToastService } from '../../../shared/components/toast/toast.service';
import { CurrencyFormatPipe } from '../../../shared/pipes/currency-format.pipe';

@Component({
  selector: 'app-invoice-list',
  standalone: true,
  imports: [PageHeaderComponent, DataTableComponent, LoadingSpinnerComponent],
  templateUrl: './invoice-list.component.html',
  styleUrls: ['./invoice-list.component.scss'],
})
export class InvoiceListComponent implements OnInit {
  private readonly invoicesService = inject(InvoicesService);
  private readonly router = inject(Router);
  private readonly confirmDialog = inject(ConfirmDialogService);
  private readonly toast = inject(ToastService);

  readonly loading = signal(true);
  readonly invoices = signal<Invoice[]>([]);

  readonly columns: TableColumn[] = [
    { key: 'invoiceNumber', label: 'Invoice #', sortable: true },
    { key: 'clientName', label: 'Client', sortable: true },
    { key: 'date', label: 'Date', sortable: true },
    { key: 'dueDate', label: 'Due Date', sortable: true },
    { key: 'total', label: 'Total', sortable: true },
    { key: 'status', label: 'Status', sortable: true },
  ];

  readonly actions: TableAction[] = [
    { label: 'Edit', action: 'edit', class: 'btn-outline' },
    { label: 'Delete', action: 'delete', class: 'btn-danger' },
  ];

  ngOnInit(): void {
    this.loadInvoices();
  }

  loadInvoices(): void {
    this.invoicesService.getAll().subscribe({
      next: (invoices) => {
        const mapped = invoices.map((inv) => ({
          ...inv,
          clientName: inv.client ? `${inv.client.firstName} ${inv.client.lastName}` : '-',
        }));
        this.invoices.set(mapped as Invoice[]);
        this.loading.set(false);
      },
      error: () => this.loading.set(false),
    });
  }

  onAdd(): void {
    this.router.navigate(['/invoices/new']);
  }

  onRowClick(row: Record<string, unknown>): void {
    this.router.navigate(['/invoices', row['id']]);
  }

  async onAction(event: { action: string; row: Record<string, unknown> }): Promise<void> {
    if (event.action === 'edit') {
      this.router.navigate(['/invoices', event.row['id'], 'edit']);
    } else if (event.action === 'delete') {
      const confirmed = await this.confirmDialog.confirm({
        title: 'Delete Invoice',
        message: `Are you sure you want to delete invoice ${event.row['invoiceNumber']}? This action cannot be undone.`,
        confirmText: 'Delete',
        confirmClass: 'btn-danger',
      });
      if (confirmed) {
        this.invoicesService.delete(event.row['id'] as number).subscribe({
          next: () => {
            this.toast.success('Invoice deleted successfully');
            this.loadInvoices();
          },
        });
      }
    }
  }
}
