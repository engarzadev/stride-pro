import { Component, inject, OnInit, signal } from '@angular/core';
import { ActivatedRoute, Router, RouterLink } from '@angular/router';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { MatTableModule } from '@angular/material/table';
import { Client } from '../../../core/models';
import { ConfirmDialogService } from '../../../shared/components/confirm-dialog/confirm-dialog.component';
import { LoadingSpinnerComponent } from '../../../shared/components/loading-spinner/loading-spinner.component';
import { ToastService } from '../../../shared/components/toast/toast.service';
import { DateFormatPipe } from '../../../shared/pipes/date-format.pipe';
import { ClientsService } from '../clients.service';

@Component({
  selector: 'app-client-detail',
  standalone: true,
  imports: [RouterLink, LoadingSpinnerComponent, DateFormatPipe, MatCardModule, MatButtonModule, MatTableModule],
  templateUrl: './client-detail.component.html',
  styleUrls: ['./client-detail.component.scss'],
})
export class ClientDetailComponent implements OnInit {
  private readonly route = inject(ActivatedRoute);
  private readonly router = inject(Router);
  private readonly clientsService = inject(ClientsService);
  private readonly confirmDialog = inject(ConfirmDialogService);
  private readonly toast = inject(ToastService);

  readonly loading = signal(true);
  readonly client = signal<Client | null>(null);

  ngOnInit(): void {
    const id = this.route.snapshot.paramMap.get('id') ?? '';
    this.clientsService.getById(id).subscribe({
      next: (client) => {
        this.client.set(client);
        this.loading.set(false);
      },
      error: () => {
        this.loading.set(false);
        this.router.navigate(['/clients']);
      },
    });
  }

  async onDelete(): Promise<void> {
    const c = this.client();
    if (!c) return;

    const confirmed = await this.confirmDialog.confirm({
      title: 'Delete Client',
      message: `Are you sure you want to delete ${c.firstName} ${c.lastName}?`,
      confirmText: 'Delete',
      confirmClass: 'btn-danger',
    });

    if (confirmed) {
      this.clientsService.delete(c.id).subscribe({
        next: () => {
          this.toast.success('Client deleted successfully');
          this.router.navigate(['/clients']);
        },
      });
    }
  }
}
