import { Component, inject, OnInit, signal, computed } from '@angular/core';
import { Router } from '@angular/router';
import { forkJoin } from 'rxjs';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { toSignal } from '@angular/core/rxjs-interop';
import { Appointment, CareLog, Client, Horse, Invoice, Reminder } from '../../core/models';
import { ApiService } from '../../core/services/api.service';
import { AuthService } from '../../core/services/auth.service';
import { HorsesService } from '../horses/horses.service';
import { LoadingSpinnerComponent } from '../../shared/components/loading-spinner/loading-spinner.component';
import { DateFormatPipe } from '../../shared/pipes/date-format.pipe';

interface UpcomingReminder {
  horseId: string;
  horseName: string;
  title: string;
  dueDate: string;
  category: string;
  urgency: 'overdue' | 'soon' | 'upcoming';
}

interface RecentCareEntry {
  horseName: string;
  category: string;
  date: string;
  notes: string;
}

const AUTO_INTERVALS: Record<string, { days: number; label: string }> = {
  deworming:   { days: 180, label: 'Deworming due' },
  vaccination: { days: 365, label: 'Vaccination due' },
  dental:      { days: 180, label: 'Dental due' },
  farrier:     { days:  42, label: 'Farrier due' },
};

@Component({
  selector: 'app-dashboard',
  standalone: true,
  imports: [LoadingSpinnerComponent, DateFormatPipe, MatCardModule, MatButtonModule, MatIconModule],
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.scss'],
})
export class DashboardComponent implements OnInit {
  private readonly api = inject(ApiService);
  private readonly router = inject(Router);
  private readonly authService = inject(AuthService);
  private readonly horsesService = inject(HorsesService);
  private readonly currentUser = toSignal(this.authService.currentUser$);

  readonly isOwner = computed(() => this.currentUser()?.role === 'owner');

  readonly loading = signal(true);

  // Professional data
  readonly clients = signal<Client[]>([]);
  readonly horses = signal<Horse[]>([]);
  readonly appointments = signal<Appointment[]>([]);
  readonly invoices = signal<Invoice[]>([]);

  // Owner data
  readonly upcomingReminders = signal<UpcomingReminder[]>([]);
  readonly recentCare = signal<RecentCareEntry[]>([]);

  ngOnInit(): void {
    if (this.isOwner()) {
      this.loadOwnerData();
    } else {
      this.loadProfessionalData();
    }
  }

  private loadProfessionalData(): void {
    forkJoin({
      clients: this.api.get<Client[]>('/clients'),
      horses: this.api.get<Horse[]>('/horses'),
      appointments: this.api.get<Appointment[]>('/appointments'),
      invoices: this.api.get<Invoice[]>('/invoices'),
    }).subscribe({
      next: (data) => {
        this.clients.set(data.clients || []);
        this.horses.set(data.horses || []);
        this.appointments.set(data.appointments || []);
        this.invoices.set(data.invoices || []);
        this.loading.set(false);
      },
      error: () => this.loading.set(false),
    });
  }

  private loadOwnerData(): void {
    this.horsesService.getAll().subscribe({
      next: (horses) => {
        this.horses.set(horses || []);
        if (horses.length === 0) {
          this.loading.set(false);
          return;
        }
        this.loadOwnerHorseData(horses);
      },
      error: () => this.loading.set(false),
    });
  }

  private loadOwnerHorseData(horses: Horse[]): void {
    const careForks = horses.map(h => this.horsesService.getCareLogs(h.id));
    const reminderForks = horses.map(h => this.horsesService.getReminders(h.id));

    forkJoin([...careForks, ...reminderForks]).subscribe({
      next: (results) => {
        const careResults = results.slice(0, horses.length) as CareLog[][];
        const reminderResults = results.slice(horses.length) as Reminder[][];
        const today = new Date();
        today.setHours(0, 0, 0, 0);

        // Build recent care entries
        const allCare: RecentCareEntry[] = [];
        for (let i = 0; i < horses.length; i++) {
          for (const log of (careResults[i] || [])) {
            allCare.push({
              horseName: horses[i].name,
              category: log.category,
              date: log.date,
              notes: log.notes,
            });
          }
        }
        allCare.sort((a, b) => b.date.localeCompare(a.date));
        this.recentCare.set(allCare.slice(0, 8));

        // Build upcoming reminders (manual + auto-generated)
        const allReminders: UpcomingReminder[] = [];
        for (let i = 0; i < horses.length; i++) {
          const horse = horses[i];
          const reminders = reminderResults[i] || [];
          const careLogs = careResults[i] || [];

          // Manual active reminders
          for (const r of reminders) {
            if (r.source === 'manual' && !r.isComplete) {
              allReminders.push({
                horseId: horse.id,
                horseName: horse.name,
                title: r.title,
                dueDate: r.dueDate,
                category: r.category,
                urgency: this.urgency(r.dueDate, today),
              });
            }
          }

          // Auto-generated reminders from care log intervals
          for (const [category, cfg] of Object.entries(AUTO_INTERVALS)) {
            const stored = reminders.find(r => r.source === 'auto' && r.category === category && !r.isComplete);
            if (stored) {
              allReminders.push({
                horseId: horse.id,
                horseName: horse.name,
                title: stored.title,
                dueDate: String(stored.dueDate).substring(0, 10),
                category,
                urgency: this.urgency(stored.dueDate, today),
              });
              continue;
            }

            const latest = careLogs
              .filter(l => l.category === category)
              .sort((a, b) => b.date.localeCompare(a.date))[0];
            if (!latest) continue;

            const dueDate = this.addDays(latest.date, cfg.days);
            const dismissed = reminders.some(
              r => r.source === 'auto' && r.category === category
                && String(r.dueDate).substring(0, 10) === dueDate && r.isComplete
            );
            if (dismissed) continue;

            allReminders.push({
              horseId: horse.id,
              horseName: horse.name,
              title: cfg.label,
              dueDate,
              category,
              urgency: this.urgency(dueDate, today),
            });
          }
        }

        allReminders.sort((a, b) => a.dueDate.localeCompare(b.dueDate));
        this.upcomingReminders.set(allReminders);
        this.loading.set(false);
      },
      error: () => this.loading.set(false),
    });
  }

  // === Professional helpers ===
  get upcomingAppointments(): Appointment[] {
    const now = new Date().toISOString();
    return this.appointments()
      .filter((a) => a.date >= now || a.status === 'scheduled')
      .slice(0, 5);
  }

  get pendingInvoices(): Invoice[] {
    return this.invoices().filter(
      (i) => i.status === 'pending' || i.status === 'overdue',
    );
  }

  // === Owner helpers ===
  get overdueReminders(): UpcomingReminder[] {
    return this.upcomingReminders().filter(r => r.urgency === 'overdue');
  }

  get soonReminders(): UpcomingReminder[] {
    return this.upcomingReminders().filter(r => r.urgency === 'soon');
  }

  categoryLabel(value: string): string {
    const labels: Record<string, string> = {
      farrier: 'Farrier', vet: 'Vet', dental: 'Dental', deworming: 'Deworming',
      vaccination: 'Vaccination', diet: 'Diet', fitting: 'Fitting',
      ride: 'Ride', training: 'Training', other: 'Other',
    };
    return labels[value] ?? value;
  }

  navigateTo(path: string): void {
    this.router.navigate([path]);
  }

  private urgency(dueDate: string, today: Date): 'overdue' | 'soon' | 'upcoming' {
    const [y, m, d] = String(dueDate).substring(0, 10).split('-').map(Number);
    const due = new Date(y, m - 1, d);
    const diff = (due.getTime() - today.getTime()) / 86_400_000;
    if (diff < 0) return 'overdue';
    if (diff <= 14) return 'soon';
    return 'upcoming';
  }

  private addDays(dateStr: string, days: number): string {
    const [y, m, d] = String(dateStr).substring(0, 10).split('-').map(Number);
    const dt = new Date(y, m - 1, d + days);
    return `${dt.getFullYear()}-${String(dt.getMonth() + 1).padStart(2, '0')}-${String(dt.getDate()).padStart(2, '0')}`;
  }
}
