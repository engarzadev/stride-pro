import { Component, inject, OnInit, signal } from '@angular/core';
import { Router } from '@angular/router';
import { BarnsService } from '../barns.service';
import { Barn } from '../../../core/models';
import { PageHeaderComponent } from '../../../shared/components/page-header/page-header.component';
import { DataTableComponent, TableColumn, TableAction } from '../../../shared/components/data-table/data-table.component';
import { LoadingSpinnerComponent } from '../../../shared/components/loading-spinner/loading-spinner.component';
import { ConfirmDialogService } from '../../../shared/components/confirm-dialog/confirm-dialog.component';
import { ToastService } from '../../../shared/components/toast/toast.service';

@Component({
  selector: 'app-barn-list',
  standalone: true,
  imports: [PageHeaderComponent, DataTableComponent, LoadingSpinnerComponent],
  templateUrl: './barn-list.component.html',
  styleUrls: ['./barn-list.component.scss'],
})
export class BarnListComponent implements OnInit {
  private readonly barnsService = inject(BarnsService);
  private readonly router = inject(Router);
  private readonly confirmDialog = inject(ConfirmDialogService);
  private readonly toast = inject(ToastService);

  readonly loading = signal(true);
  readonly barns = signal<Barn[]>([]);

  readonly columns: TableColumn[] = [
    { key: 'name', label: 'Name', sortable: true },
    { key: 'address', label: 'Address', sortable: true },
    { key: 'phone', label: 'Phone' },
    { key: 'email', label: 'Email' },
  ];

  readonly actions: TableAction[] = [
    { label: 'Edit', action: 'edit', class: 'btn-outline' },
    { label: 'Delete', action: 'delete', class: 'btn-danger' },
  ];

  ngOnInit(): void {
    this.loadBarns();
  }

  loadBarns(): void {
    this.barnsService.getAll().subscribe({
      next: (barns) => {
        this.barns.set(barns);
        this.loading.set(false);
      },
      error: () => this.loading.set(false),
    });
  }

  onAdd(): void {
    this.router.navigate(['/barns/new']);
  }

  onRowClick(row: Record<string, unknown>): void {
    this.router.navigate(['/barns', row['id']]);
  }

  async onAction(event: { action: string; row: Record<string, unknown> }): Promise<void> {
    if (event.action === 'edit') {
      this.router.navigate(['/barns', event.row['id'], 'edit']);
    } else if (event.action === 'delete') {
      const confirmed = await this.confirmDialog.confirm({
        title: 'Delete Barn',
        message: `Are you sure you want to delete ${event.row['name']}? This action cannot be undone.`,
        confirmText: 'Delete',
        confirmClass: 'btn-danger',
      });
      if (confirmed) {
        this.barnsService.delete(event.row['id'] as number).subscribe({
          next: () => {
            this.toast.success('Barn deleted successfully');
            this.loadBarns();
          },
        });
      }
    }
  }
}
