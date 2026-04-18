import { Component } from '@angular/core';
import { MatTabsModule } from '@angular/material/tabs';
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
export class SettingsComponent {}
