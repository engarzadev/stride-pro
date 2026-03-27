import { Component, inject, OnInit, signal } from '@angular/core';
import { ActivatedRoute, Router, RouterLink } from '@angular/router';
import { SessionsService } from '../sessions.service';
import { Session } from '../../../core/models';
import { LoadingSpinnerComponent } from '../../../shared/components/loading-spinner/loading-spinner.component';
import { ConfirmDialogService } from '../../../shared/components/confirm-dialog/confirm-dialog.component';
import { ToastService } from '../../../shared/components/toast/toast.service';
import { DateFormatPipe } from '../../../shared/pipes/date-format.pipe';

@Component({
  selector: 'app-session-detail',
  standalone: true,
  imports: [RouterLink, LoadingSpinnerComponent, DateFormatPipe],
  templateUrl: './session-detail.component.html',
  styleUrls: ['./session-detail.component.scss'],
})
export class SessionDetailComponent implements OnInit {
  private readonly route = inject(ActivatedRoute);
  private readonly router = inject(Router);
  private readonly sessionsService = inject(SessionsService);
  private readonly confirmDialog = inject(ConfirmDialogService);
  private readonly toast = inject(ToastService);

  readonly loading = signal(true);
  readonly session = signal<Session | null>(null);

  ngOnInit(): void {
    const id = Number(this.route.snapshot.paramMap.get('id'));
    this.sessionsService.getById(id).subscribe({
      next: (session) => {
        this.session.set(session);
        this.loading.set(false);
      },
      error: () => {
        this.loading.set(false);
        this.router.navigate(['/sessions']);
      },
    });
  }

  async onDelete(): Promise<void> {
    const s = this.session();
    if (!s) return;

    const confirmed = await this.confirmDialog.confirm({
      title: 'Delete Session',
      message: 'Are you sure you want to delete this session?',
      confirmText: 'Delete',
      confirmClass: 'btn-danger',
    });

    if (confirmed) {
      this.sessionsService.delete(s.id).subscribe({
        next: () => {
          this.toast.success('Session deleted successfully');
          this.router.navigate(['/sessions']);
        },
      });
    }
  }
}
