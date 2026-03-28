import { Component, inject, OnInit, signal } from '@angular/core';
import { ActivatedRoute, Router, RouterLink } from '@angular/router';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { MatTableModule } from '@angular/material/table';
import { BarnsService } from '../barns.service';
import { Barn } from '../../../core/models';
import { LoadingSpinnerComponent } from '../../../shared/components/loading-spinner/loading-spinner.component';
import { ConfirmDialogService } from '../../../shared/components/confirm-dialog/confirm-dialog.component';
import { ToastService } from '../../../shared/components/toast/toast.service';
import { DateFormatPipe } from '../../../shared/pipes/date-format.pipe';

@Component({
  selector: 'app-barn-detail',
  standalone: true,
  imports: [RouterLink, LoadingSpinnerComponent, DateFormatPipe, MatCardModule, MatButtonModule, MatTableModule],
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
