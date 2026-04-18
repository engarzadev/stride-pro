import { Injectable, inject } from '@angular/core';
import { Observable } from 'rxjs';
import { ApiService } from '../../core/services/api.service';
import { BusinessSettings, ServiceItem } from '../../core/models';

@Injectable({ providedIn: 'root' })
export class SettingsService {
  private readonly api = inject(ApiService);

  getBusinessSettings(): Observable<BusinessSettings> {
    return this.api.get<BusinessSettings>('/settings/business');
  }

  saveBusinessSettings(settings: BusinessSettings): Observable<BusinessSettings> {
    return this.api.put<BusinessSettings>('/settings/business', settings);
  }

  getServiceItems(): Observable<ServiceItem[]> {
    return this.api.get<ServiceItem[]>('/settings/service-items');
  }

  createServiceItem(item: { name: string; defaultPrice: number }): Observable<ServiceItem> {
    return this.api.post<ServiceItem>('/settings/service-items', item);
  }

  updateServiceItem(id: string, item: { name: string; defaultPrice: number }): Observable<ServiceItem> {
    return this.api.put<ServiceItem>(`/settings/service-items/${id}`, item);
  }

  deleteServiceItem(id: string): Observable<void> {
    return this.api.delete<void>(`/settings/service-items/${id}`);
  }
}
