import { DatePipe } from '@angular/common';
import { Component, computed, inject, OnInit, signal } from '@angular/core';
import { Router } from '@angular/router';
import { MatCardModule } from '@angular/material/card';
import { Invoice } from '../../../core/models';
import { ConfirmDialogService } from '../../../shared/components/confirm-dialog/confirm-dialog.component';
import {
  DataTableComponent,
  TableAction,
  TableColumn,
} from '../../../shared/components/data-table/data-table.component';
import { LoadingSpinnerComponent } from '../../../shared/components/loading-spinner/loading-spinner.component';
import { PageHeaderComponent } from '../../../shared/components/page-header/page-header.component';
import { ToastService } from '../../../shared/components/toast/toast.service';
import { InvoicesService } from '../invoices.service';

@Component({
  selector: 'app-invoice-list',
  standalone: true,
  imports: [PageHeaderComponent, DataTableComponent, LoadingSpinnerComponent, MatCardModule],
  templateUrl: './invoice-list.component.html',
  styleUrls: ['./invoice-list.component.scss'],
})
export class InvoiceListComponent implements OnInit {
  private readonly invoicesService = inject(InvoicesService);
  private readonly router = inject(Router);
  private readonly confirmDialog = inject(ConfirmDialogService);
  private readonly toast = inject(ToastService);
  private readonly datePipe = new DatePipe('en-US');

  readonly loading = signal(true);
  readonly invoices = signal<Invoice[]>([]);

  readonly tableInvoices = computed(() =>
    this.invoices().map((inv) => ({
      ...inv,
      invoiceNumber: `#${inv.id.slice(0, 5).toUpperCase()}`,
      clientName: inv.client
        ? `${inv.client.firstName} ${inv.client.lastName}`
        : '',
      date: this.datePipe.transform(inv.createdAt, 'MMMM dd, yyyy') ?? '',
      dueDate: this.datePipe.transform(inv.dueDate, 'MMMM dd, yyyy') ?? '',
      total: `$${inv.total.toFixed(2)}`,
    })),
  );

  readonly columns: TableColumn[] = [
    { key: 'invoiceNumber', label: 'Invoice #', sortable: true },
    { key: 'clientName', label: 'Client', sortable: true },
    { key: 'date', label: 'Date', sortable: true },
    { key: 'dueDate', label: 'Due Date', sortable: true },
    { key: 'total', label: 'Total', sortable: true },
    {
      key: 'status',
      label: 'Status',
      sortable: true,
      type: 'badge',
      capitalize: true,
      badgeMap: {
        draft: 'secondary',
        sent: 'primary',
        paid: 'success',
        overdue: 'danger',
      },
    },
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
        this.invoices.set(invoices);
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

  async onAction(event: {
    action: string;
    row: Record<string, unknown>;
  }): Promise<void> {
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
        this.invoicesService.delete(event.row['id'] as string).subscribe({
          next: () => {
            this.toast.success('Invoice deleted successfully');
            this.loadInvoices();
          },
        });
      }
    }
  }
}
