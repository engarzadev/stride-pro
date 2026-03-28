import { Component, inject, OnInit, signal } from '@angular/core';
import { TitleCasePipe } from '@angular/common';
import { ActivatedRoute, Router, RouterLink } from '@angular/router';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { MatTableModule } from '@angular/material/table';
import { InvoicesService } from '../invoices.service';
import { Invoice } from '../../../core/models';
import { LoadingSpinnerComponent } from '../../../shared/components/loading-spinner/loading-spinner.component';
import { ConfirmDialogService } from '../../../shared/components/confirm-dialog/confirm-dialog.component';
import { ToastService } from '../../../shared/components/toast/toast.service';
import { DateFormatPipe } from '../../../shared/pipes/date-format.pipe';
import { CurrencyFormatPipe } from '../../../shared/pipes/currency-format.pipe';

@Component({
  selector: 'app-invoice-detail',
  standalone: true,
  imports: [RouterLink, LoadingSpinnerComponent, DateFormatPipe, CurrencyFormatPipe, TitleCasePipe, MatCardModule, MatButtonModule, MatTableModule],
  templateUrl: './invoice-detail.component.html',
  styleUrls: ['./invoice-detail.component.scss'],
})
export class InvoiceDetailComponent implements OnInit {
  private readonly route = inject(ActivatedRoute);
  private readonly router = inject(Router);
  private readonly invoicesService = inject(InvoicesService);
  private readonly confirmDialog = inject(ConfirmDialogService);
  private readonly toast = inject(ToastService);

  readonly loading = signal(true);
  readonly invoice = signal<Invoice | null>(null);

  ngOnInit(): void {
    const id = this.route.snapshot.paramMap.get('id')!;
    this.invoicesService.getById(id).subscribe({
      next: (invoice) => {
        this.invoice.set(invoice);
        this.loading.set(false);
      },
      error: () => {
        this.loading.set(false);
        this.router.navigate(['/invoices']);
      },
    });
  }

  getStatusClass(status: string): string {
    switch (status) {
      case 'paid': return 'badge-success';
      case 'sent': return 'badge-primary';
      case 'overdue': return 'badge-danger';
      case 'draft': return 'badge-secondary';
      default: return 'badge-secondary';
    }
  }

  async onDelete(): Promise<void> {
    const inv = this.invoice();
    if (!inv) return;

    const confirmed = await this.confirmDialog.confirm({
      title: 'Delete Invoice',
      message: `Are you sure you want to delete invoice ${inv.invoiceNumber}?`,
      confirmText: 'Delete',
      confirmClass: 'btn-danger',
    });

    if (confirmed) {
      this.invoicesService.delete(inv.id).subscribe({
        next: () => {
          this.toast.success('Invoice deleted successfully');
          this.router.navigate(['/invoices']);
        },
      });
    }
  }
}
