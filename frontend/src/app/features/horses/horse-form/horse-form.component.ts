import { Component, inject, OnInit, signal } from '@angular/core';
import { ActivatedRoute, Router, RouterLink } from '@angular/router';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { HorsesService } from '../horses.service';
import { ClientsService } from '../../clients/clients.service';
import { BarnsService } from '../../barns/barns.service';
import { Client, Barn } from '../../../core/models';
import { ToastService } from '../../../shared/components/toast/toast.service';
import { LoadingSpinnerComponent } from '../../../shared/components/loading-spinner/loading-spinner.component';
import { BreedAutocompleteComponent } from '../../../shared/components/breed-autocomplete/breed-autocomplete.component';

@Component({
  selector: 'app-horse-form',
  standalone: true,
  imports: [ReactiveFormsModule, RouterLink, LoadingSpinnerComponent, BreedAutocompleteComponent],
  templateUrl: './horse-form.component.html',
  styleUrls: ['./horse-form.component.scss'],
})
export class HorseFormComponent implements OnInit {
  private readonly fb = inject(FormBuilder);
  private readonly route = inject(ActivatedRoute);
  private readonly router = inject(Router);
  private readonly horsesService = inject(HorsesService);
  private readonly clientsService = inject(ClientsService);
  private readonly barnsService = inject(BarnsService);
  private readonly toast = inject(ToastService);

  readonly loading = signal(false);
  readonly saving = signal(false);
  readonly isEdit = signal(false);
  readonly clients = signal<Client[]>([]);
  readonly barns = signal<Barn[]>([]);
  private horseId = 0;

  readonly form = this.fb.nonNullable.group({
    name: ['', [Validators.required]],
    breed: [''],
    age: [0],
    gender: [''],
    color: [''],
    weight: [0],
    notes: [''],
    clientId: [0, [Validators.required, Validators.min(1)]],
    barnId: [0],
  });

  ngOnInit(): void {
    this.clientsService.getAll().subscribe((c) => this.clients.set(c));
    this.barnsService.getAll().subscribe((b) => this.barns.set(b));

    const id = this.route.snapshot.paramMap.get('id');
    if (id) {
      this.isEdit.set(true);
      this.horseId = Number(id);
      this.loading.set(true);
      this.horsesService.getById(this.horseId).subscribe({
        next: (horse) => {
          this.form.patchValue({
            name: horse.name,
            breed: horse.breed,
            age: horse.age,
            gender: horse.gender,
            color: horse.color,
            weight: horse.weight,
            notes: horse.notes,
            clientId: horse.clientId,
            barnId: horse.barnId,
          });
          this.loading.set(false);
        },
        error: () => {
          this.loading.set(false);
          this.router.navigate(['/horses']);
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
      ? this.horsesService.update(this.horseId, data)
      : this.horsesService.create(data);

    request$.subscribe({
      next: () => {
        this.toast.success(this.isEdit() ? 'Horse updated successfully' : 'Horse created successfully');
        this.router.navigate(['/horses']);
      },
      error: () => this.saving.set(false),
    });
  }
}
