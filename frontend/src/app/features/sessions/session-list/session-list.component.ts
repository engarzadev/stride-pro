import { Component, computed, inject, OnInit, signal } from '@angular/core';
import { MatCardModule } from '@angular/material/card';
import { MatIconModule } from '@angular/material/icon';
import { Router, RouterLink } from '@angular/router';
import { Session } from '../../../core/models';
import { SubscriptionService } from '../../../core/services/subscription.service';
import { ConfirmDialogService } from '../../../shared/components/confirm-dialog/confirm-dialog.component';
import {
  DataTableComponent,
  FilterConfig,
  MobileCardConfig,
  TableAction,
  TableColumn,
} from '../../../shared/components/data-table/data-table.component';
import { LoadingSpinnerComponent } from '../../../shared/components/loading-spinner/loading-spinner.component';
import { PageHeaderComponent } from '../../../shared/components/page-header/page-header.component';
import { ToastService } from '../../../shared/components/toast/toast.service';
import { SessionsService } from '../sessions.service';

@Component({
  selector: 'app-session-list',
  standalone: true,
  imports: [
    PageHeaderComponent,
    DataTableComponent,
    LoadingSpinnerComponent,
    MatCardModule,
    MatIconModule,
    RouterLink
  ],
  templateUrl: './session-list.component.html',
  styleUrls: ['./session-list.component.scss'],
})
export class SessionListComponent implements OnInit {
  private readonly sessionsService = inject(SessionsService);
  private readonly router = inject(Router);
  private readonly confirmDialog = inject(ConfirmDialogService);
  private readonly toast = inject(ToastService);
  private readonly subscriptionService = inject(SubscriptionService);

  readonly loading = signal(true);
  readonly sessions = signal<Session[]>([]);
  readonly canManageSessions = computed(() => this.subscriptionService.hasFeature('session_notes'));

  readonly tableSessions = computed(() =>
    this.sessions().map((s) => ({
      ...s,
      clientName: s.appointment?.client
        ? `${s.appointment.client.firstName} ${s.appointment.client.lastName}`
        : '',
      horseName: s.appointment?.horse?.name ?? '',
    })),
  );

  readonly columns: TableColumn[] = [
    { key: 'createdAt', label: 'Date', sortable: true, type: 'date' },
    { key: 'type', label: 'Type', sortable: true, capitalize: true },
    { key: 'horseName', label: 'Horse', sortable: true },
    { key: 'clientName', label: 'Client', sortable: true },
  ];

  readonly actions: TableAction[] = [
    { label: 'Edit', action: 'edit', class: 'btn-outline' },
    { label: 'Delete', action: 'delete', class: 'btn-danger' },
  ];

  readonly mobileCard: MobileCardConfig = {
    titleKey: 'createdAt',
    subtitleKey: 'horseName',
  };

  readonly filterConfig: FilterConfig[] = [
    {
      key: 'type',
      label: 'Type',
      options: [
        { value: 'massage', label: 'Massage' },
        { value: 'chiropractic', label: 'Chiropractic' },
        { value: 'acupuncture', label: 'Acupuncture' },
        { value: 'rehabilitation', label: 'Rehabilitation' },
        { value: 'evaluation', label: 'Evaluation' },
        { value: 'treatment', label: 'Treatment' },
        { value: 'pemf', label: 'PEMF' },
        { value: 'other', label: 'Other' },
      ],
    },
  ];

  ngOnInit(): void {
    this.loadSessions();
  }

  loadSessions(): void {
    this.sessionsService.getAll().subscribe({
      next: (sessions) => {
        this.sessions.set(sessions);
        this.loading.set(false);
      },
      error: () => this.loading.set(false),
    });
  }

  onAdd(): void {
    this.router.navigate(['/sessions/new']);
  }

  onRowClick(row: Record<string, unknown>): void {
    this.router.navigate(['/sessions', row['id']]);
  }

  async onAction(event: {
    action: string;
    row: Record<string, unknown>;
  }): Promise<void> {
    if (event.action === 'edit') {
      this.router.navigate(['/sessions', event.row['id'], 'edit']);
    } else if (event.action === 'delete') {
      const confirmed = await this.confirmDialog.confirm({
        title: 'Delete Session',
        message:
          'Are you sure you want to delete this session? This action cannot be undone.',
        confirmText: 'Delete',
        confirmClass: 'btn-danger',
      });
      if (confirmed) {
        this.sessionsService.delete(event.row['id'] as string).subscribe({
          next: () => {
            this.toast.success('Session deleted successfully');
            this.loadSessions();
          },
        });
      }
    }
  }
}
