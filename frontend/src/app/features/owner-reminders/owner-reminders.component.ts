import { Component, inject, OnInit, signal, computed } from '@angular/core';
import { MatSelectModule } from '@angular/material/select';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { RouterLink } from '@angular/router';
import { Horse } from '../../core/models';
import { HorsesService } from '../horses/horses.service';
import { HorseRemindersComponent } from '../horses/horse-reminders/horse-reminders.component';
import { LoadingSpinnerComponent } from '../../shared/components/loading-spinner/loading-spinner.component';
import { SubscriptionService } from '../../core/services/subscription.service';

@Component({
  selector: 'app-owner-reminders',
  standalone: true,
  imports: [
    MatSelectModule,
    MatFormFieldModule,
    MatIconModule,
    RouterLink,
    HorseRemindersComponent,
    LoadingSpinnerComponent,
  ],
  templateUrl: './owner-reminders.component.html',
  styleUrls: ['./owner-reminders.component.scss'],
})
export class OwnerRemindersComponent implements OnInit {
  private readonly horsesService = inject(HorsesService);
  private readonly subscriptionService = inject(SubscriptionService);

  readonly loading = signal(true);
  readonly canUseCareLog = computed(() => this.subscriptionService.hasFeature('care_logs'));
  readonly horses = signal<Horse[]>([]);
  readonly selectedHorseId = signal<string | null>(null);

  readonly selectedHorse = computed(() => {
    const id = this.selectedHorseId();
    return this.horses().find(h => h.id === id) ?? null;
  });

  ngOnInit(): void {
    this.horsesService.getAll().subscribe({
      next: (horses) => {
        this.horses.set(horses);
        if (horses.length > 0) {
          this.selectedHorseId.set(horses[0].id);
        }
        this.loading.set(false);
      },
      error: () => this.loading.set(false),
    });
  }

  onHorseChange(horseId: string): void {
    this.selectedHorseId.set(horseId);
  }
}
