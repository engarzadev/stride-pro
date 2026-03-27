import { Component, inject, OnInit, signal } from '@angular/core';
import { Router } from '@angular/router';
import { ClientsService } from '../clients.service';
import { Client } from '../../../core/models';
import { PageHeaderComponent } from '../../../shared/components/page-header/page-header.component';
import { DataTableComponent, TableColumn, TableAction } from '../../../shared/components/data-table/data-table.component';
import { LoadingSpinnerComponent } from '../../../shared/components/loading-spinner/loading-spinner.component';
import { ConfirmDialogService } from '../../../shared/components/confirm-dialog/confirm-dialog.component';
import { ToastService } from '../../../shared/components/toast/toast.service';

@Component({
  selector: 'app-client-list',
  standalone: true,
  imports: [PageHeaderComponent, DataTableComponent, LoadingSpinnerComponent],
  templateUrl: './client-list.component.html',
  styleUrls: ['./client-list.component.scss'],
})
export class ClientListComponent implements OnInit {
  private readonly clientsService = inject(ClientsService);
  private readonly router = inject(Router);
  private readonly confirmDialog = inject(ConfirmDialogService);
  private readonly toast = inject(ToastService);

  readonly loading = signal(true);
  readonly clients = signal<Client[]>([]);

  readonly columns: TableColumn[] = [
    { key: 'firstName', label: 'First Name', sortable: true },
    { key: 'lastName', label: 'Last Name', sortable: true },
    { key: 'email', label: 'Email', sortable: true },
    { key: 'phone', label: 'Phone' },
  ];

  readonly actions: TableAction[] = [
    { label: 'Edit', action: 'edit', class: 'btn-outline' },
    { label: 'Delete', action: 'delete', class: 'btn-danger' },
  ];

  ngOnInit(): void {
    this.loadClients();
  }

  loadClients(): void {
    this.clientsService.getAll().subscribe({
      next: (clients) => {
        this.clients.set(clients);
        this.loading.set(false);
      },
      error: () => this.loading.set(false),
    });
  }

  onAdd(): void {
    this.router.navigate(['/clients/new']);
  }

  onRowClick(row: Record<string, unknown>): void {
    this.router.navigate(['/clients', row['id']]);
  }

  async onAction(event: { action: string; row: Record<string, unknown> }): Promise<void> {
    if (event.action === 'edit') {
      this.router.navigate(['/clients', event.row['id'], 'edit']);
    } else if (event.action === 'delete') {
      const confirmed = await this.confirmDialog.confirm({
        title: 'Delete Client',
        message: `Are you sure you want to delete ${event.row['firstName']} ${event.row['lastName']}? This action cannot be undone.`,
        confirmText: 'Delete',
        confirmClass: 'btn-danger',
      });
      if (confirmed) {
        this.clientsService.delete(event.row['id'] as number).subscribe({
          next: () => {
            this.toast.success('Client deleted successfully');
            this.loadClients();
          },
        });
      }
    }
  }
}
