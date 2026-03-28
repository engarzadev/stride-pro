import { Injectable, inject, signal } from '@angular/core';
import { Observable, of, tap } from 'rxjs';
import { map } from 'rxjs/operators';
import { SubscriptionPlan, SubscriptionResponse } from '../models';
import { ApiService } from './api.service';

@Injectable({ providedIn: 'root' })
export class SubscriptionService {
  private readonly api = inject(ApiService);

  private readonly _plan = signal<SubscriptionPlan | null>(null);
  readonly plan = this._plan.asReadonly();

  private loaded = false;

  load(): Observable<SubscriptionPlan | null> {
    if (this.loaded) {
      return of(this._plan());
    }
    return this.api.get<SubscriptionResponse>('/subscription').pipe(
      tap((resp) => {
        this._plan.set(resp.plan);
        this.loaded = true;
      }),
      map((resp) => resp.plan),
    );
  }

  hasFeature(feature: string): boolean {
    return this._plan()?.features?.includes(feature) ?? false;
  }
}
