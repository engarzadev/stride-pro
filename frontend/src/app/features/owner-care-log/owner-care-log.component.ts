import { Component, inject, OnInit, signal, computed } from '@angular/core';
import { MatSelectModule } from '@angular/material/select';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { Horse } from '../../core/models';
import { HorsesService } from '../horses/horses.service';
import { CareLogComponent } from '../horses/care-log/care-log.component';
import { LoadingSpinnerComponent } from '../../shared/components/loading-spinner/loading-spinner.component';

@Component({
  selector: 'app-owner-care-log',
  standalone: true,
  imports: [
    MatSelectModule,
    MatFormFieldModule,
    MatIconModule,
    CareLogComponent,
    LoadingSpinnerComponent,
  ],
  templateUrl: './owner-care-log.component.html',
  styleUrls: ['./owner-care-log.component.scss'],
})
export class OwnerCareLogComponent implements OnInit {
  private readonly horsesService = inject(HorsesService);

  readonly loading = signal(true);
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
