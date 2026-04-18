import { Component, computed, inject, OnInit, signal } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { ActivatedRoute, Router } from '@angular/router';
import { Barn } from '../../../core/models';
import { ConfirmDialogService } from '../../../shared/components/confirm-dialog/confirm-dialog.component';
import {
  DataTableComponent,
  MobileCardConfig,
  TableColumn,
} from '../../../shared/components/data-table/data-table.component';
import { DetailPageComponent } from '../../../shared/components/detail-page/detail-page.component';
import { LoadingSpinnerComponent } from '../../../shared/components/loading-spinner/loading-spinner.component';
import { ToastService } from '../../../shared/components/toast/toast.service';
import { DateFormatPipe } from '../../../shared/pipes/date-format.pipe';
import { BarnsService } from '../barns.service';

@Component({
  selector: 'app-barn-detail',
  standalone: true,
  imports: [
    LoadingSpinnerComponent,
    DateFormatPipe,
    MatCardModule,
    MatButtonModule,
    DataTableComponent,
    DetailPageComponent,
  ],
  templateUrl: './barn-detail.component.html',
  styleUrls: ['./barn-detail.component.scss'],
})
export class BarnDetailComponent implements OnInit {
  private readonly route = inject(ActivatedRoute);
  private readonly router = inject(Router);
  private readonly barnsService = inject(BarnsService);
  private readonly confirmDialog = inject(ConfirmDialogService);
  private readonly toast = inject(ToastService);

  readonly loading = signal(true);
  readonly barn = signal<Barn | null>(null);

  readonly horseColumns: TableColumn[] = [
    { key: 'name', label: 'Name', sortable: true },
    { key: 'breed', label: 'Breed' },
    { key: 'ownerName', label: 'Owner' },
  ];

  readonly horseMobileCard: MobileCardConfig = {
    titleKey: 'name',
    subtitleKey: 'breed',
  };

  readonly horseRows = computed(() => {
    const b = this.barn();
    if (!b?.horses) return [];
    return b.horses.map((h) => ({
      ...h,
      ownerName: h.client ? `${h.client.firstName} ${h.client.lastName}` : '-',
    }));
  });

  onHorseClick(row: Record<string, unknown>): void {
    this.router.navigate(['/horses', row['id']]);
  }

  ngOnInit(): void {
    const id = this.route.snapshot.paramMap.get('id')!;
    this.barnsService.getById(id).subscribe({
      next: (barn) => {
        this.barn.set(barn);
        this.loading.set(false);
      },
      error: () => {
        this.loading.set(false);
        this.router.navigate(['/barns']);
      },
    });
  }

  async onDelete(): Promise<void> {
    const b = this.barn();
    if (!b) return;

    const confirmed = await this.confirmDialog.confirm({
      title: 'Delete Barn',
      message: `Are you sure you want to delete ${b.name}?`,
      confirmText: 'Delete',
      confirmClass: 'btn-danger',
    });

    if (confirmed) {
      this.barnsService.delete(b.id).subscribe({
        next: () => {
          this.toast.success('Barn deleted successfully');
          this.router.navigate(['/barns']);
        },
      });
    }
  }
}
