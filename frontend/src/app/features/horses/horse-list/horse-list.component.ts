import { Component, computed, inject, OnInit, signal } from '@angular/core';
import { Router } from '@angular/router';
import { MatCardModule } from '@angular/material/card';
import { Horse } from '../../../core/models';
import { ConfirmDialogService } from '../../../shared/components/confirm-dialog/confirm-dialog.component';
import {
  DataTableComponent,
  TableAction,
  TableColumn,
} from '../../../shared/components/data-table/data-table.component';
import { LoadingSpinnerComponent } from '../../../shared/components/loading-spinner/loading-spinner.component';
import { PageHeaderComponent } from '../../../shared/components/page-header/page-header.component';
import { ToastService } from '../../../shared/components/toast/toast.service';
import { HorsesService } from '../horses.service';

@Component({
  selector: 'app-horse-list',
  standalone: true,
  imports: [PageHeaderComponent, DataTableComponent, LoadingSpinnerComponent, MatCardModule],
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

  readonly tableHorses = computed(() =>
    this.horses().map((h) => ({
      ...h,
      clientName: h.client ? `${h.client.firstName} ${h.client.lastName}` : '',
      barnName: h.barn?.name ?? '',
    })),
  );

  readonly columns: TableColumn[] = [
    { key: 'name', label: 'Name', sortable: true },
    { key: 'breed', label: 'Breed', sortable: true },
    { key: 'clientName', label: 'Client', sortable: true },
    { key: 'barnName', label: 'Barn', sortable: true },
    { key: 'age', label: 'Age', sortable: true },
    { key: 'gender', label: 'Gender', capitalize: true },
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

  async onAction(event: {
    action: string;
    row: Record<string, unknown>;
  }): Promise<void> {
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
        this.horsesService.delete(event.row['id'] as string).subscribe({
          next: () => {
            this.toast.success('Horse deleted successfully');
            this.loadHorses();
          },
        });
      }
    }
  }
}
