import { Component, computed, inject } from '@angular/core';
import { toSignal } from '@angular/core/rxjs-interop';
import { MatTabsModule } from '@angular/material/tabs';
import { AuthService } from '../../core/services/auth.service';
import { BusinessSettingsComponent } from './business-settings/business-settings.component';
import { ServiceCatalogComponent } from './service-catalog/service-catalog.component';
import { GeneralSettingsComponent } from './general-settings/general-settings.component';

@Component({
  selector: 'app-settings',
  standalone: true,
  imports: [MatTabsModule, BusinessSettingsComponent, ServiceCatalogComponent, GeneralSettingsComponent],
  templateUrl: './settings.component.html',
  styleUrls: ['./settings.component.scss'],
})
export class SettingsComponent {
  private readonly authService = inject(AuthService);
  private readonly currentUser = toSignal(this.authService.currentUser$);

  readonly isProfessional = computed(() => this.currentUser()?.role !== 'owner');
}
