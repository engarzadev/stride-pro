# Page Templates

Canonical templates for list, form, and detail pages. Copy these when scaffolding a new feature.

Each feature lives at `src/app/features/[entity]/` and follows this file structure:

```
features/[entity]/
├── [entity]-list/
│   ├── [entity]-list.component.ts
│   └── [entity]-list.component.html
├── [entity]-form/
│   ├── [entity]-form.component.ts
│   └── [entity]-form.component.html
├── [entity]-detail/
│   ├── [entity]-detail.component.ts
│   └── [entity]-detail.component.html
├── [entity].service.ts
└── [entity].routes.ts
```

---

## Routes

**`[entity].routes.ts`**

```typescript
import { Routes } from '@angular/router';
import { EntityListComponent } from './entity-list/entity-list.component';
import { EntityFormComponent } from './entity-form/entity-form.component';
import { EntityDetailComponent } from './entity-detail/entity-detail.component';

export const entityRoutes: Routes = [
  { path: '', component: EntityListComponent },
  { path: 'new', component: EntityFormComponent },
  { path: ':id', component: EntityDetailComponent },
  { path: ':id/edit', component: EntityFormComponent },
];
```

Register in `app.routes.ts` with lazy loading:

```typescript
{
  path: 'entities',
  loadChildren: () => import('./features/entity/entity.routes').then(m => m.entityRoutes),
  canActivate: [authGuard],
},
```

---

## List Page

**`[entity]-list.component.ts`**

```typescript
import { Component, inject, OnInit, signal } from '@angular/core';
import { Router } from '@angular/router';
import { MatCardModule } from '@angular/material/card';
import { Entity } from '../../../core/models';
import { EntityService } from '../entity.service';
import { LoadingSpinnerComponent } from '../../../shared/components/loading-spinner/loading-spinner.component';
import { PageHeaderComponent } from '../../../shared/components/page-header/page-header.component';
import { DataTableComponent, MobileCardConfig, TableColumn, TableAction } from '../../../shared/components/data-table/data-table.component';
import { ConfirmDialogService } from '../../../shared/components/confirm-dialog/confirm-dialog.component';
import { ToastService } from '../../../shared/components/toast/toast.service';

@Component({
  selector: 'app-entity-list',
  standalone: true,
  imports: [LoadingSpinnerComponent, PageHeaderComponent, DataTableComponent, MatCardModule],
  templateUrl: './entity-list.component.html',
  styleUrls: ['./entity-list.component.scss'],
})
export class EntityListComponent implements OnInit {
  private readonly router = inject(Router);
  private readonly entityService = inject(EntityService);
  private readonly confirmDialog = inject(ConfirmDialogService);
  private readonly toast = inject(ToastService);

  readonly loading = signal(true);
  readonly entities = signal<Entity[]>([]);

  readonly columns: TableColumn[] = [
    { key: 'name', label: 'Name', sortable: true },
    // { key: 'nested.field', label: 'Nested', sortable: true },   // dot-path for nested values
    // { key: 'status', label: 'Status', type: 'badge', badgeMap: { active: 'primary', inactive: 'danger' } },
    // { key: 'date', label: 'Date', type: 'date' },
    // { key: 'amount', label: 'Amount', type: 'currency' },
  ];

  readonly actions: TableAction[] = [
    { label: 'Edit', action: 'edit', class: 'btn-outline' },
    { label: 'Delete', action: 'delete', class: 'btn-danger' },
  ];

  // titleKey: primary card heading on mobile; subtitleKey: secondary line
  readonly mobileCard: MobileCardConfig = { titleKey: 'name', subtitleKey: 'someOtherField' };

  ngOnInit(): void {
    this.entityService.getAll().subscribe({
      next: (entities) => {
        this.entities.set(entities);
        this.loading.set(false);
      },
      error: () => this.loading.set(false),
    });
  }

  onAdd(): void {
    this.router.navigate(['/entities/new']);
  }

  onRowClick(row: Record<string, unknown>): void {
    this.router.navigate(['/entities', row['id']]);
  }

  async onAction(event: { action: string; row: Record<string, unknown> }): Promise<void> {
    if (event.action === 'edit') {
      this.router.navigate(['/entities', event.row['id'], 'edit']);
    }
    if (event.action === 'delete') {
      const confirmed = await this.confirmDialog.confirm({
        title: 'Delete Entity',
        message: `Are you sure you want to delete this entity?`,
        confirmText: 'Delete',
        confirmClass: 'btn-danger',
      });
      if (confirmed) {
        this.entityService.delete(event.row['id'] as string).subscribe({
          next: () => {
            this.toast.success('Entity deleted successfully');
            this.entities.update(list => list.filter(e => e.id !== event.row['id']));
          },
        });
      }
    }
  }
}
```

**`[entity]-list.component.html`**

```html
<app-page-header
  title="Entities"
  buttonText="Add Entity"
  (buttonClick)="onAdd()" />

@if (loading()) {
  <app-loading-spinner />
} @else {
  <mat-card>
    <mat-card-content>
      <app-data-table
        [columns]="columns"
        [data]="$any(entities())"
        [actions]="actions"
        [mobileCard]="mobileCard"
        (rowClick)="onRowClick($event)"
        (actionClick)="onAction($event)" />
    </mat-card-content>
  </mat-card>
}
```

---

## Form Page

Uses `<app-form-page>` from `shared/components/form-page/`. The wrapper renders the back link, page title, and Save/Cancel actions card. You only write the form content cards.

**`[entity]-form.component.ts`**

```typescript
import { Component, inject, OnInit, signal } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { EntityService } from '../entity.service';
import { FormPageComponent } from '../../../shared/components/form-page/form-page.component';
import { LoadingSpinnerComponent } from '../../../shared/components/loading-spinner/loading-spinner.component';
import { ToastService } from '../../../shared/components/toast/toast.service';

@Component({
  selector: 'app-entity-form',
  standalone: true,
  imports: [
    ReactiveFormsModule, FormPageComponent, LoadingSpinnerComponent,
    MatCardModule, MatButtonModule, MatIconModule, MatFormFieldModule, MatInputModule,
  ],
  templateUrl: './entity-form.component.html',
  styleUrls: ['./entity-form.component.scss'],
})
export class EntityFormComponent implements OnInit {
  private readonly route = inject(ActivatedRoute);
  private readonly router = inject(Router);
  private readonly fb = inject(FormBuilder);
  private readonly entityService = inject(EntityService);
  private readonly toast = inject(ToastService);

  readonly loading = signal(false);
  readonly saving = signal(false);
  readonly isEdit = signal(false);
  private entityId = '';

  readonly form = this.fb.nonNullable.group({
    name: ['', [Validators.required]],
    // email: ['', [Validators.email]],
    // phone: [''],
    // notes: [''],
  });

  ngOnInit(): void {
    const id = this.route.snapshot.paramMap.get('id');
    if (id) {
      this.isEdit.set(true);
      this.entityId = id;
      this.loading.set(true);
      this.entityService.getById(id).subscribe({
        next: (entity) => {
          this.form.patchValue(entity);
          this.loading.set(false);
        },
        error: () => {
          this.loading.set(false);
          this.router.navigate(['/entities']);
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
      ? this.entityService.update(this.entityId, data)
      : this.entityService.create(data);

    request$.subscribe({
      next: () => {
        this.toast.success(this.isEdit() ? 'Entity updated' : 'Entity created');
        this.router.navigate(['/entities']);
      },
      error: () => this.saving.set(false),
    });
  }
}
```

**`[entity]-form.component.html`**

```html
@if (loading()) {
  <app-loading-spinner />
} @else {
  <app-form-page
    [title]="isEdit() ? 'Edit Entity' : 'New Entity'"
    backRoute="/entities"
    backLabel="Back to Entities"
    [saving]="saving()"
    [isEdit]="isEdit()"
    createIcon="add"
    (save)="onSubmit()">
    <form [formGroup]="form">
      <div class="form-main">

        <mat-card class="form-card">
          <mat-card-content>
            <div class="card-section-title">Entity Information</div>
            <div class="form-row">
              <!-- 2-column grid on desktop, 1-column on mobile -->
              <mat-form-field appearance="outline">
                <mat-label>Name</mat-label>
                <input matInput type="text" formControlName="name" />
                @if (form.controls.name.errors?.['required']) {
                  <mat-error>Name is required.</mat-error>
                }
              </mat-form-field>
            </div>
          </mat-card-content>
        </mat-card>

        <!-- Notes card (add when entity has notes) -->
        <mat-card class="form-card">
          <mat-card-content>
            <div class="card-section-title">Notes</div>
            <mat-form-field appearance="outline" class="full-width">
              <mat-label>Notes (optional)</mat-label>
              <textarea matInput formControlName="notes" rows="4"></textarea>
            </mat-form-field>
          </mat-card-content>
        </mat-card>

      </div>
    </form>
  </app-form-page>
}
```

### `createIcon` reference

| Entity type | Icon |
|---|---|
| Generic | `add` |
| Client/person | `person_add` |
| Appointment/event | `event` |
| Invoice | `send` |

### Field reference

| Field type | Template snippet |
|---|---|
| Text | `<input matInput type="text" formControlName="x" />` |
| Email | `<input matInput type="email" formControlName="x" />` |
| Number | `<input matInput type="number" formControlName="x" />` |
| Textarea | `<textarea matInput formControlName="x" rows="4"></textarea>` |
| Select | `<mat-select formControlName="x"><mat-option [value]="v">Label</mat-option></mat-select>` |
| Date | `<input matInput [matDatepicker]="dp" formControlName="x" /><mat-datepicker-toggle /><mat-datepicker #dp />` |

---

## Detail Page

Uses `<app-detail-page>` from `shared/components/detail-page/`. The wrapper renders the back link, page title, and Edit/Delete buttons. You only write the content cards.

**`[entity]-detail.component.ts`**

```typescript
import { Component, inject, OnInit, signal } from '@angular/core';
import { ActivatedRoute, Router, RouterLink } from '@angular/router';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { Entity } from '../../../core/models';
import { EntityService } from '../entity.service';
import { DetailPageComponent } from '../../../shared/components/detail-page/detail-page.component';
import { LoadingSpinnerComponent } from '../../../shared/components/loading-spinner/loading-spinner.component';
import { ConfirmDialogService } from '../../../shared/components/confirm-dialog/confirm-dialog.component';
import { ToastService } from '../../../shared/components/toast/toast.service';
import { DateFormatPipe } from '../../../shared/pipes/date-format.pipe';

@Component({
  selector: 'app-entity-detail',
  standalone: true,
  imports: [RouterLink, LoadingSpinnerComponent, DateFormatPipe, MatCardModule, MatButtonModule, DetailPageComponent],
  templateUrl: './entity-detail.component.html',
  styleUrls: ['./entity-detail.component.scss'],
})
export class EntityDetailComponent implements OnInit {
  private readonly route = inject(ActivatedRoute);
  private readonly router = inject(Router);
  private readonly entityService = inject(EntityService);
  private readonly confirmDialog = inject(ConfirmDialogService);
  private readonly toast = inject(ToastService);

  readonly loading = signal(true);
  readonly entity = signal<Entity | null>(null);

  ngOnInit(): void {
    const id = this.route.snapshot.paramMap.get('id')!;
    this.entityService.getById(id).subscribe({
      next: (entity) => {
        this.entity.set(entity);
        this.loading.set(false);
      },
      error: () => {
        this.loading.set(false);
        this.router.navigate(['/entities']);
      },
    });
  }

  async onDelete(): Promise<void> {
    const e = this.entity();
    if (!e) return;

    const confirmed = await this.confirmDialog.confirm({
      title: 'Delete Entity',
      message: `Are you sure you want to delete this entity?`,
      confirmText: 'Delete',
      confirmClass: 'btn-danger',
    });

    if (confirmed) {
      this.entityService.delete(e.id).subscribe({
        next: () => {
          this.toast.success('Entity deleted successfully');
          this.router.navigate(['/entities']);
        },
      });
    }
  }
}
```

**`[entity]-detail.component.html`**

```html
@if (loading()) {
  <app-loading-spinner />
} @else {
  @if (entity(); as e) {
    <app-detail-page
      [title]="e.name"
      backRoute="/entities"
      backLabel="Back to Entities"
      [editRoute]="['/entities', e.id, 'edit']"
      (deleteClick)="onDelete()">

      <!-- Primary info card -->
      <mat-card class="mb-6">
        <mat-card-header>
          <mat-card-title>Entity Information</mat-card-title>
        </mat-card-header>
        <mat-card-content>
          <div class="detail-grid">
            <!-- Each detail-field = one label/value pair -->
            <div class="detail-field">
              <label>Field Label</label>
              <p>{{ e.someField || '-' }}</p>
            </div>
            <div class="detail-field">
              <label>Date</label>
              <p>{{ e.createdAt | dateFormat }}</p>
            </div>
            <!-- Linked entity -->
            <div class="detail-field">
              <label>Related</label>
              @if (e.related) {
                <p><a [routerLink]="['/related', e.relatedId]">{{ e.related.name }}</a></p>
              } @else {
                <p>-</p>
              }
            </div>
          </div>
          @if (e.notes) {
            <div class="detail-field mt-4">
              <label>Notes</label>
              <p>{{ e.notes }}</p>
            </div>
          }
        </mat-card-content>
      </mat-card>

      <!-- Related entity list (add DataTableComponent to TS imports if needed) -->
      @if (e.children && e.children.length > 0) {
        <mat-card>
          <mat-card-header>
            <mat-card-title>Related Items</mat-card-title>
          </mat-card-header>
          <mat-card-content class="p-0">
            <app-data-table
              [columns]="childColumns"
              [data]="$any(e.children)"
              [mobileCard]="childMobileCard"
              (rowClick)="onChildClick($event)" />
          </mat-card-content>
        </mat-card>
      }

    </app-detail-page>
  }
}
```

> **Related entity table:** Add `DataTableComponent, TableColumn, MobileCardConfig` to TS imports. Define `childColumns`, `childMobileCard`, and an `onChildClick()` method — same pattern as client-detail and barn-detail.

---

## Service

**`[entity].service.ts`**

```typescript
import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { Entity } from '../../core/models';
import { environment } from '../../../environments/environment';

@Injectable({ providedIn: 'root' })
export class EntityService {
  private readonly http = inject(HttpClient);
  private readonly base = `${environment.apiUrl}/entities`;

  getAll(): Observable<Entity[]> {
    return this.http.get<Entity[]>(this.base);
  }

  getById(id: string): Observable<Entity> {
    return this.http.get<Entity>(`${this.base}/${id}`);
  }

  create(data: Partial<Entity>): Observable<Entity> {
    return this.http.post<Entity>(this.base, data);
  }

  update(id: string, data: Partial<Entity>): Observable<Entity> {
    return this.http.patch<Entity>(`${this.base}/${id}`, data);
  }

  delete(id: string): Observable<void> {
    return this.http.delete<void>(`${this.base}/${id}`);
  }
}
```
