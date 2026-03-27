import { Component, inject, OnInit, signal } from '@angular/core';
import { ActivatedRoute, Router, RouterLink } from '@angular/router';
import { HorsesService } from '../horses.service';
import { Horse } from '../../../core/models';
import { LoadingSpinnerComponent } from '../../../shared/components/loading-spinner/loading-spinner.component';
import { ConfirmDialogService } from '../../../shared/components/confirm-dialog/confirm-dialog.component';
import { ToastService } from '../../../shared/components/toast/toast.service';
import { DateFormatPipe } from '../../../shared/pipes/date-format.pipe';

@Component({
  selector: 'app-horse-detail',
  standalone: true,
  imports: [RouterLink, LoadingSpinnerComponent, DateFormatPipe],
  templateUrl: './horse-detail.component.html',
  styleUrls: ['./horse-detail.component.scss'],
})
export class HorseDetailComponent implements OnInit {
  private readonly route = inject(ActivatedRoute);
  private readonly router = inject(Router);
  private readonly horsesService = inject(HorsesService);
  private readonly confirmDialog = inject(ConfirmDialogService);
  private readonly toast = inject(ToastService);

  readonly loading = signal(true);
  readonly horse = signal<Horse | null>(null);

  ngOnInit(): void {
    const id = Number(this.route.snapshot.paramMap.get('id'));
    this.horsesService.getById(id).subscribe({
      next: (horse) => {
        this.horse.set(horse);
        this.loading.set(false);
      },
      error: () => {
        this.loading.set(false);
        this.router.navigate(['/horses']);
      },
    });
  }

  async onDelete(): Promise<void> {
    const h = this.horse();
    if (!h) return;

    const confirmed = await this.confirmDialog.confirm({
      title: 'Delete Horse',
      message: `Are you sure you want to delete ${h.name}?`,
      confirmText: 'Delete',
      confirmClass: 'btn-danger',
    });

    if (confirmed) {
      this.horsesService.delete(h.id).subscribe({
        next: () => {
          this.toast.success('Horse deleted successfully');
          this.router.navigate(['/horses']);
        },
      });
    }
  }
}
