import { Component, computed, inject, input, OnInit, signal } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatDatepickerModule } from '@angular/material/datepicker';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';
import { CareLog, Reminder } from '../../../core/models';
import { ToastService } from '../../../shared/components/toast/toast.service';
import { DateFormatPipe } from '../../../shared/pipes/date-format.pipe';
import { CARE_LOG_CATEGORIES } from '../care-log/care-log.component';
import { HorsesService } from '../horses.service';
import { ActivatedRoute } from '@angular/router';

const AUTO_INTERVALS: Record<string, { days: number; label: string }> = {
  deworming:   { days: 180, label: 'Deworming due' },
  vaccination: { days: 365, label: 'Vaccination due' },
  dental:      { days: 180, label: 'Dental due' },
  farrier:     { days:  42, label: 'Farrier due' },
};

export interface DisplayReminder {
  id: string;
  title: string;
  dueDate: string;
  category: string;
  source: 'manual' | 'auto';
  isVirtual: boolean;
  urgency: 'overdue' | 'soon' | 'upcoming';
}

@Component({
  selector: 'app-horse-reminders',
  standalone: true,
  imports: [
    ReactiveFormsModule,
    DateFormatPipe,
    MatButtonModule,
    MatDatepickerModule,
    MatFormFieldModule,
    MatIconModule,
    MatInputModule,
    MatSelectModule,
  ],
  templateUrl: './horse-reminders.component.html',
  styleUrls: ['./horse-reminders.component.scss'],
})
export class HorseRemindersComponent implements OnInit {
  readonly horseId = input.required<string>();

  private readonly horsesService = inject(HorsesService);
  private readonly toast = inject(ToastService);
  private readonly fb = inject(FormBuilder);
  private readonly activatedRoute = inject(ActivatedRoute);

  private readonly reminders = signal<Reminder[]>([]);
  private readonly careLogs = signal<CareLog[]>([]);

  readonly loading = signal(true);
  readonly showForm = signal(false);
  readonly showCompleted = signal(false);
  readonly saving = signal(false);
  readonly editingId = signal<string | null>(null);
  readonly editingItemId = signal<string | null>(null);
  readonly editingSource = signal<'manual' | 'auto'>('manual');

  readonly categories = CARE_LOG_CATEGORIES;

  readonly form = this.fb.group({
    title:    ['', Validators.required],
    dueDate:  ['', Validators.required],
    category: [''],
  });

  readonly activeReminders = computed(() => this.buildList(false));
  readonly completedReminders = computed(() => this.buildList(true));

  ngOnInit(): void {
    const openForm = this.activatedRoute.snapshot.paramMap.get('showForm');
    if (openForm) {
      this.onAdd();
    }

    this.load();
  }

  private load(): void {
    this.loading.set(true);
    let remindersLoaded = false;
    let logsLoaded = false;

    const tryDone = () => {
      if (remindersLoaded && logsLoaded) this.loading.set(false);
    };

    this.horsesService.getReminders(this.horseId()).subscribe({
      next: (r) => { this.reminders.set(r); remindersLoaded = true; tryDone(); },
      error: () => { remindersLoaded = true; tryDone(); },
    });

    this.horsesService.getCareLogs(this.horseId()).subscribe({
      next: (l) => { this.careLogs.set(l); logsLoaded = true; tryDone(); },
      error: () => { logsLoaded = true; tryDone(); },
    });
  }

  private buildList(complete: boolean): DisplayReminder[] {
    const today = new Date();
    today.setHours(0, 0, 0, 0);
    const items: DisplayReminder[] = [];

    for (const r of this.reminders()) {
      if (r.source === 'manual' && r.isComplete === complete) {
        items.push({ id: r.id, title: r.title, dueDate: r.dueDate, category: r.category, source: 'manual', isVirtual: false, urgency: this.urgency(r.dueDate, today) });
      }
      if (complete && r.source === 'auto' && r.isComplete) {
        items.push({ id: r.id, title: r.title, dueDate: r.dueDate, category: r.category, source: 'auto', isVirtual: false, urgency: this.urgency(r.dueDate, today) });
      }
    }

    if (!complete) {
      for (const [category, cfg] of Object.entries(AUTO_INTERVALS)) {
        // Prefer a stored (edited) active auto reminder over the computed one
        const stored = this.reminders().find(
          r => r.source === 'auto' && r.category === category && !r.isComplete
        );
        if (stored) {
          items.push({ id: stored.id, title: stored.title, dueDate: String(stored.dueDate).substring(0, 10), category, source: 'auto', isVirtual: false, urgency: this.urgency(stored.dueDate, today) });
          continue;
        }

        const latest = this.careLogs()
          .filter(l => l.category === category)
          .sort((a, b) => b.date.localeCompare(a.date))[0];
        if (!latest) continue;

        const dueDate = this.addDays(latest.date, cfg.days);
        const dismissed = this.reminders().some(
          r => r.source === 'auto' && r.category === category
            && String(r.dueDate).substring(0, 10) === dueDate && r.isComplete
        );
        if (dismissed) continue;

        items.push({ id: `auto-${category}`, title: cfg.label, dueDate, category, source: 'auto', isVirtual: true, urgency: this.urgency(dueDate, today) });
      }
    }

    const editingItemId = this.editingItemId();
    return items
      .filter(i => i.id !== editingItemId)
      .sort((a, b) =>
        complete ? b.dueDate.localeCompare(a.dueDate) : a.dueDate.localeCompare(b.dueDate)
      );
  }

  onAdd(): void {
    this.editingId.set(null);
    this.editingSource.set('manual');
    this.form.reset();
    this.showForm.set(true);
  }

  onEdit(item: DisplayReminder): void {
    this.editingId.set(item.isVirtual ? null : item.id);
    this.editingItemId.set(item.id);
    this.editingSource.set(item.source);
    const [y, m, d] = item.dueDate.split('-').map(Number);
    this.form.setValue({
      title:    item.title,
      dueDate:  new Date(y, m - 1, d) as unknown as string,
      category: item.category,
    });
    this.showForm.set(true);
  }

  onCancelForm(): void {
    this.showForm.set(false);
    this.editingId.set(null);
    this.editingItemId.set(null);
    this.editingSource.set('manual');
    this.form.reset();
  }

  onSave(): void {
    if (this.form.invalid) return;
    this.saving.set(true);

    const { title, dueDate, category } = this.form.value;
    const d = dueDate as unknown as Date;
    const dueDateStr = d instanceof Date ? this.formatDate(d) : dueDate!;
    const editId = this.editingId();

    const req$ = editId
      ? this.horsesService.putReminder(editId, { title: title!, dueDate: dueDateStr, category: category ?? '' })
      : this.horsesService.createReminder(this.horseId(), { title: title!, dueDate: dueDateStr, category: category ?? '', source: this.editingSource() });

    req$.subscribe({
      next: (r) => {
        this.reminders.update(list =>
          editId ? list.map(x => x.id === r.id ? r : x) : [...list, r]
        );
        this.toast.success(editId ? 'Reminder updated' : 'Reminder added');
        this.onCancelForm();
        this.saving.set(false);
      },
      error: () => {
        this.toast.error(editId ? 'Failed to update reminder' : 'Failed to add reminder');
        this.saving.set(false);
      },
    });
  }

  onMarkComplete(item: DisplayReminder): void {
    if (item.isVirtual) {
      this.horsesService.createReminder(this.horseId(), {
        title: item.title, dueDate: item.dueDate, category: item.category, source: 'auto', isComplete: true,
      }).subscribe({
        next: (r) => { this.reminders.update(list => [...list, r]); },
        error: () => this.toast.error('Failed to update reminder'),
      });
    } else {
      this.horsesService.patchReminder(item.id, { isComplete: true }).subscribe({
        next: (r) => { this.reminders.update(list => list.map(x => x.id === r.id ? r : x)); },
        error: () => this.toast.error('Failed to update reminder'),
      });
    }
  }

  onReopen(item: DisplayReminder): void {
    this.horsesService.patchReminder(item.id, { isComplete: false }).subscribe({
      next: (r) => { this.reminders.update(list => list.map(x => x.id === r.id ? r : x)); },
      error: () => this.toast.error('Failed to update reminder'),
    });
  }

  onDelete(item: DisplayReminder): void {
    if (item.isVirtual) {
      this.horsesService.createReminder(this.horseId(), {
        title: item.title, dueDate: item.dueDate, category: item.category, source: 'auto', isComplete: true,
      }).subscribe({
        next: (r) => { this.reminders.update(list => [...list, r]); },
        error: () => this.toast.error('Failed to delete reminder'),
      });
    } else {
      this.horsesService.deleteReminder(item.id).subscribe({
        next: () => { this.reminders.update(list => list.filter(x => x.id !== item.id)); },
        error: () => this.toast.error('Failed to delete reminder'),
      });
    }
  }

  categoryLabel(value: string): string {
    return this.categories.find(c => c.value === value)?.label ?? '';
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
    return this.formatDate(new Date(y, m - 1, d + days));
  }

  private formatDate(d: Date): string {
    return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`;
  }
}
