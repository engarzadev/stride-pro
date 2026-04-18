import { Component, inject } from '@angular/core';
import { MatCardModule } from '@angular/material/card';
import { MatSlideToggleModule } from '@angular/material/slide-toggle';
import { ThemeService } from '../../../core/services/theme.service';

@Component({
  selector: 'app-general-settings',
  standalone: true,
  imports: [MatCardModule, MatSlideToggleModule],
  templateUrl: './general-settings.component.html',
  styleUrls: ['./general-settings.component.scss'],
})
export class GeneralSettingsComponent {
  private readonly themeService = inject(ThemeService);
  readonly isDark = this.themeService.isDark;

  onToggle(): void {
    this.themeService.toggle();
  }
}
