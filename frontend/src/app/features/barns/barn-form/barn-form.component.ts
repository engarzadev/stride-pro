import { Component, inject, OnInit, signal } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { FormPageComponent } from '../../../shared/components/form-page/form-page.component';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatCardModule } from '@angular/material/card';
import { SubscriptionService } from '../../../core/services/subscription.service';
import { BarnsService } from '../barns.service';
import { HorsesService } from '../../horses/horses.service';
import { ToastService } from '../../../shared/components/toast/toast.service';
import { LoadingSpinnerComponent } from '../../../shared/components/loading-spinner/loading-spinner.component';
import { UpgradeFieldPromptComponent } from '../../../shared/components/upgrade-field-prompt/upgrade-field-prompt.component';
import { HorseMultiselectAutocompleteComponent } from '../../../shared/components/horse-multiselect-autocomplete/horse-multiselect-autocomplete.component';
import { Horse } from '../../../core/models';
import { Observable, forkJoin, of } from 'rxjs';

@Component({
  selector: 'app-barn-form',
  standalone: true,
  imports: [ReactiveFormsModule, FormPageComponent, LoadingSpinnerComponent, HorseMultiselectAutocompleteComponent, UpgradeFieldPromptComponent, MatFormFieldModule, MatInputModule, MatButtonModule, MatIconModule, MatCardModule],
  templateUrl: './barn-form.component.html',
  styleUrls: ['./barn-form.component.scss'],
})
export class BarnFormComponent implements OnInit {
  private readonly fb = inject(FormBuilder);
  private readonly route = inject(ActivatedRoute);
  private readonly router = inject(Router);
  private readonly barnsService = inject(BarnsService);
  private readonly horsesService = inject(HorsesService);
  private readonly toast = inject(ToastService);
  private readonly subscriptionService = inject(SubscriptionService);

  readonly loading = signal(false);
  readonly saving = signal(false);
  readonly isEdit = signal(false);
  readonly canManageBarns = signal(false);
  readonly allHorses = signal<Horse[]>([]);
  readonly selectedHorse = signal<Horse[]>([]);
  private originalHorse: Horse[] = [];
  private barnId = '';

  readonly form = this.fb.nonNullable.group({
    name: ['', [Validators.required]],
    contactName: [''],
    address: [''],
    phone: [''],
    email: ['', [Validators.email]],
    notes: [''],
  });

  ngOnInit(): void {
    const id = this.route.snapshot.paramMap.get('id');

    this.subscriptionService.load().subscribe(() => {
      const has = this.subscriptionService.hasFeature('barn_management');
      this.canManageBarns.set(has);

      if (!has && !id) {
        this.toast.error('Barn management requires a paid plan');
        this.router.navigate(['/barns']);
        return;
      }

      if (has) {
        this.horsesService.getAll().subscribe((horses) => {
          this.allHorses.set(horses);
          if (id) {
            const assigned = horses.filter((h) => h.barnId === id);
            this.selectedHorse.set(assigned);
            this.originalHorse = assigned;
          }
        });
      }
    });

    if (id) {
      this.isEdit.set(true);
      this.barnId = id;
      this.loading.set(true);
      this.barnsService.getById(this.barnId).subscribe({
        next: (barn) => {
          this.form.patchValue({
            name: barn.name,
            contactName: barn.contactName,
            address: barn.address,
            phone: barn.phone,
            email: barn.email,
            notes: barn.notes,
          });
          this.loading.set(false);
        },
        error: () => {
          this.loading.set(false);
          this.router.navigate(['/barns']);
        },
      });
    }
  }

  onSubmit(): void {
    if (this.form.invalid) {
      this.form.markAllAsTouched();
      return;
    }

    this.saving.set(true);
    const data = this.form.getRawValue();

    const request$ = this.isEdit()
      ? this.barnsService.update(this.barnId, data)
      : this.barnsService.create(data);

    request$.subscribe({
      next: (barn) => {
        this.saveHorseAssignments(barn.id).subscribe({
          next: () => {
            this.toast.success(this.isEdit() ? 'Barn updated successfully' : 'Barn created successfully');
            this.router.navigate(['/barns']);
          },
          error: () => this.saving.set(false),
        });
      },
      error: () => this.saving.set(false),
    });
  }

  private saveHorseAssignments(barnId: string): Observable<unknown> {
    const originalSet = new Set(this.originalHorse.map((h) => h.id));
    const newSet = new Set(this.selectedHorse().map((h) => h.id));

    const toAdd = this.selectedHorse().filter((h) => !originalSet.has(h.id));
    const toRemove = this.originalHorse.filter((h) => !newSet.has(h.id));

    const updates = [
      ...toAdd.map((horse) => this.horsesService.update(horse.id, { ...horse, barnId })),
      ...toRemove.map((horse) => this.horsesService.update(horse.id, { ...horse, barnId: null })),
    ];

    return updates.length > 0 ? forkJoin(updates) : of(null);
  }
}
