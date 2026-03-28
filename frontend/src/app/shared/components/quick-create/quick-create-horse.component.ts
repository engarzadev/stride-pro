import { Component, Injectable, OnInit, computed, effect, inject, signal } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { Client, Barn, Horse } from '../../../core/models';
import { HorsesService } from '../../../features/horses/horses.service';
import { ClientsService } from '../../../features/clients/clients.service';
import { BarnsService } from '../../../features/barns/barns.service';
import { BreedAutocompleteComponent } from '../breed-autocomplete/breed-autocomplete.component';
import { QuickCreateBarnService } from './quick-create-barn.component';

export interface QuickCreateHorseOptions {
  clientId?: string;
}

@Injectable({ providedIn: 'root' })
export class QuickCreateHorseService {
  readonly visible = signal(false);
  readonly options = signal<QuickCreateHorseOptions>({});
  private resolveFn?: (result: Horse | null) => void;

  open(options: QuickCreateHorseOptions = {}): Promise<Horse | null> {
    this.options.set(options);
    this.visible.set(true);
    return new Promise<Horse | null>((resolve) => {
      this.resolveFn = resolve;
    });
  }

  complete(horse: Horse): void {
    this.visible.set(false);
    this.resolveFn?.(horse);
  }

  cancel(): void {
    this.visible.set(false);
    this.resolveFn?.(null);
  }
}

@Component({
  selector: 'app-quick-create-horse',
  standalone: true,
  imports: [ReactiveFormsModule, BreedAutocompleteComponent],
  templateUrl: './quick-create-horse.component.html',
  styleUrls: ['./quick-create-modal.scss'],
})
export class QuickCreateHorseComponent implements OnInit {
  readonly service = inject(QuickCreateHorseService);
  private readonly horsesService = inject(HorsesService);
  private readonly clientsService = inject(ClientsService);
  private readonly barnsService = inject(BarnsService);
  private readonly quickCreateBarn = inject(QuickCreateBarnService);
  private readonly fb = inject(FormBuilder);

  readonly saving = signal(false);
  readonly clients = signal<Client[]>([]);
  readonly barns = signal<Barn[]>([]);
  readonly preselectedClient = computed(() => {
    const clientId = this.service.options().clientId;
    return clientId ? this.clients().find((c) => c.id === clientId) ?? null : null;
  });

  readonly form = this.fb.nonNullable.group({
    name: ['', [Validators.required]],
    breed: [''],
    age: [0],
    gender: [''],
    color: [''],
    weight: [0],
    notes: [''],
    clientId: ['', [Validators.required]],
    barnId: [null as string | null],
  });

  constructor() {
    effect(() => {
      if (this.service.visible()) {
        const clientId = this.service.options().clientId;
        if (clientId) {
          this.form.controls.clientId.setValue(clientId);
        }
      }
    });
  }

  ngOnInit(): void {
    this.clientsService.getAll().subscribe((c) => this.clients.set(c));
    this.barnsService.getAll().subscribe((b) => this.barns.set(b));
  }

  onSubmit(): void {
    if (this.form.invalid) {
      this.form.markAllAsTouched();
      return;
    }
    this.saving.set(true);
    const data = this.form.getRawValue();
    this.horsesService.create(data).subscribe({
      next: (horse) => {
        this.saving.set(false);
        this.form.reset({ name: '', breed: '', age: 0, gender: '', color: '', weight: 0, notes: '', clientId: '', barnId: null });
        this.service.complete(horse);
      },
      error: () => this.saving.set(false),
    });
  }

  async openCreateBarn(): Promise<void> {
    const barn = await this.quickCreateBarn.open();
    if (barn) {
      this.barns.update((b) => [...b, barn]);
      this.form.controls.barnId.setValue(barn.id);
    }
  }

  onCancel(): void {
    this.form.reset({ name: '', breed: '', age: 0, gender: '', color: '', weight: 0, notes: '', clientId: '', barnId: null });
    this.service.cancel();
  }
}
