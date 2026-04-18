import { Component, computed, inject, input, OnInit, signal } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { ScrollingModule } from '@angular/cdk/scrolling';
import { MatButtonModule } from '@angular/material/button';
import { MatDatepickerModule } from '@angular/material/datepicker';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';
import { RouterLink } from '@angular/router';
import { CareLog } from '../../../core/models';
import { SubscriptionService } from '../../../core/services/subscription.service';
import { ConfirmDialogService } from '../../../shared/components/confirm-dialog/confirm-dialog.component';
import { LoadingSpinnerComponent } from '../../../shared/components/loading-spinner/loading-spinner.component';
import { ToastService } from '../../../shared/components/toast/toast.service';
import { DateFormatPipe } from '../../../shared/pipes/date-format.pipe';
import { HorsesService } from '../horses.service';

export const CARE_LOG_CATEGORIES = [
  { value: 'farrier', label: 'Farrier' },
  { value: 'vet', label: 'Vet' },
  { value: 'dental', label: 'Dental' },
  { value: 'deworming', label: 'Deworming' },
  { value: 'vaccination', label: 'Vaccination' },
  { value: 'diet', label: 'Diet' },
  { value: 'fitting', label: 'Fitting' },
  { value: 'ride', label: 'Ride' },
  { value: 'training', label: 'Training' },
  { value: 'other', label: 'Other' },
];

@Component({
  selector: 'app-care-log',
  standalone: true,
  imports: [
    ReactiveFormsModule,
    RouterLink,
    DateFormatPipe,
    ScrollingModule,
    MatButtonModule,
    MatFormFieldModule,
    MatIconModule,
    MatInputModule,
    MatDatepickerModule,
    MatSelectModule,
    LoadingSpinnerComponent,
  ],
  templateUrl: './care-log.component.html',
  styleUrls: ['./care-log.component.scss'],
})
export class CareLogComponent implements OnInit {
  readonly horseId = input.required<string>();

  private readonly horsesService = inject(HorsesService);
  private readonly subscriptionService = inject(SubscriptionService);
  private readonly confirmDialog = inject(ConfirmDialogService);
  private readonly toast = inject(ToastService);
  private readonly fb = inject(FormBuilder);

  readonly loading = signal(true);
  readonly canUseCareLog = signal(false);
  readonly logs = signal<CareLog[]>([]);
  readonly showForm = signal(false);
  readonly editingId = signal<string | null>(null);
  readonly saving = signal(false);

  readonly filterCategory = signal<string>('');
  readonly sortDir = signal<'asc' | 'desc'>('desc');

  readonly filteredLogs = computed(() => {
    const cat = this.filterCategory();
    const dir = this.sortDir();
    return this.logs()
      .filter(l => !cat || l.category === cat)
      .sort((a, b) => {
        const cmp = a.date.localeCompare(b.date);
        return dir === 'desc' ? -cmp : cmp;
      });
  });

  readonly visibleLogs = computed(() =>
    this.filteredLogs().filter(l => l.id !== this.editingId())
  );

  readonly viewportHeight = computed(() =>
    Math.min(this.visibleLogs().length * 88, 440)
  );

  trackLog(_: number, log: CareLog): string { return log.id; }

  readonly categories = CARE_LOG_CATEGORIES;

  readonly form = this.fb.group({
    date: ['', Validators.required],
    category: ['', Validators.required],
    notes: [''],
  });

  ngOnInit(): void {
    this.subscriptionService.load().subscribe(() => {
      const allowed = this.subscriptionService.hasFeature('care_logs');
      this.canUseCareLog.set(allowed);
      if (allowed) {
        this.loadLogs();
      } else {
        this.loading.set(false);
      }
    });
  }

  loadLogs(): void {
    this.horsesService.getCareLogs(this.horseId()).subscribe({
      next: (logs) => {
        this.logs.set(logs);
        this.loading.set(false);
      },
      error: () => this.loading.set(false),
    });
  }

  onAdd(): void {
    this.editingId.set(null);
    this.form.reset();
    this.showForm.set(true);
  }

  onEdit(log: CareLog): void {
    this.editingId.set(log.id);
    const [y, m, d] = log.date.split('-').map(Number);
    this.form.setValue({ date: new Date(y, m - 1, d) as unknown as string, category: log.category, notes: log.notes });
    this.showForm.set(true);
  }

  onCancelForm(): void {
    this.showForm.set(false);
    this.editingId.set(null);
    this.form.reset();
  }

  onSave(): void {
    if (this.form.invalid) return;
    this.saving.set(true);

    const { date, category, notes } = this.form.value;
    const d = date as unknown as Date;
    const dateStr = d instanceof Date
      ? `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
      : date!;
    const payload = { date: dateStr, category: category!, notes: notes ?? '' };
    const editId = this.editingId();

    const req$ = editId
      ? this.horsesService.updateCareLog(editId, payload)
      : this.horsesService.createCareLog(this.horseId(), payload);

    req$.subscribe({
      next: () => {
        this.toast.success(editId ? 'Entry updated' : 'Entry added');
        this.onCancelForm();
        this.loadLogs();
        this.saving.set(false);
      },
      error: () => {
        this.toast.error('Failed to save entry');
        this.saving.set(false);
      },
    });
  }

  async onDelete(log: CareLog): Promise<void> {
    const confirmed = await this.confirmDialog.confirm({
      title: 'Delete Care Log Entry',
      message: 'Are you sure you want to delete this entry?',
      confirmText: 'Delete',
      confirmClass: 'btn-danger',
    });
    if (!confirmed) return;

    this.horsesService.deleteCareLog(log.id).subscribe({
      next: () => {
        this.toast.success('Entry deleted');
        this.loadLogs();
      },
    });
  }

  categoryLabel(value: string): string {
    return this.categories.find((c) => c.value === value)?.label ?? value;
  }
}
