import { Component, inject, OnInit, signal } from '@angular/core';
import { Router } from '@angular/router';
import { HorsesService } from '../horses.service';
import { Horse } from '../../../core/models';
import { PageHeaderComponent } from '../../../shared/components/page-header/page-header.component';
import { DataTableComponent, TableColumn, TableAction } from '../../../shared/components/data-table/data-table.component';
import { LoadingSpinnerComponent } from '../../../shared/components/loading-spinner/loading-spinner.component';
import { ConfirmDialogService } from '../../../shared/components/confirm-dialog/confirm-dialog.component';
import { ToastService } from '../../../shared/components/toast/toast.service';

@Component({
  selector: 'app-horse-list',
  standalone: true,
  imports: [PageHeaderComponent, DataTableComponent, LoadingSpinnerComponent],
  templateUrl: './horse-list.component.html',
  styleUrls: ['./horse-list.component.scss'],
})
export class HorseListComponent implements OnInit {
  private readonly horsesService = inject(HorsesService);
  private readonly router = inject(Router);
  private readonly confirmDialog = inject(ConfirmDialogService);
  private readonly toast = inject(ToastService);

  readonly loading = signal(true);
  readonly horses = signal<Horse[]>([]);

  readonly columns: TableColumn[] = [
    { key: 'name', label: 'Name', sortable: true },
    { key: 'breed', label: 'Breed', sortable: true },
    { key: 'client.firstName', label: 'Client', sortable: true },
    { key: 'age', label: 'Age', sortable: true },
    { key: 'gender', label: 'Gender' },
  ];

  readonly actions: TableAction[] = [
    { label: 'Edit', action: 'edit', class: 'btn-outline' },
    { label: 'Delete', action: 'delete', class: 'btn-danger' },
  ];

  ngOnInit(): void {
    this.loadHorses();
  }

  loadHorses(): void {
    this.horsesService.getAll().subscribe({
      next: (horses) => {
        this.horses.set(horses);
        this.loading.set(false);
      },
      error: () => this.loading.set(false),
    });
  }

  onAdd(): void {
    this.router.navigate(['/horses/new']);
  }

  onRowClick(row: Record<string, unknown>): void {
    this.router.navigate(['/horses', row['id']]);
  }

  async onAction(event: { action: string; row: Record<string, unknown> }): Promise<void> {
    if (event.action === 'edit') {
      this.router.navigate(['/horses', event.row['id'], 'edit']);
    } else if (event.action === 'delete') {
      const confirmed = await this.confirmDialog.confirm({
        title: 'Delete Horse',
        message: `Are you sure you want to delete ${event.row['name']}? This action cannot be undone.`,
        confirmText: 'Delete',
        confirmClass: 'btn-danger',
      });
      if (confirmed) {
        this.horsesService.delete(event.row['id'] as number).subscribe({
          next: () => {
            this.toast.success('Horse deleted successfully');
            this.loadHorses();
          },
        });
      }
    }
  }
}
