import { Component, inject, signal } from '@angular/core';
import { Router, RouterOutlet, NavigationEnd } from '@angular/router';
import { filter } from 'rxjs';
import { HeaderComponent } from './shared/components/header/header.component';
import { SidebarComponent } from './shared/components/sidebar/sidebar.component';
import { ToastComponent } from './shared/components/toast/toast.component';
import { ConfirmDialogComponent } from './shared/components/confirm-dialog/confirm-dialog.component';
import { QuickCreateClientComponent } from './shared/components/quick-create/quick-create-client.component';
import { QuickCreateBarnComponent } from './shared/components/quick-create/quick-create-barn.component';
import { QuickCreateHorseComponent } from './shared/components/quick-create/quick-create-horse.component';
import { QuickCreateAppointmentComponent } from './shared/components/quick-create/quick-create-appointment.component';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [RouterOutlet, HeaderComponent, SidebarComponent, ToastComponent, ConfirmDialogComponent, QuickCreateClientComponent, QuickCreateBarnComponent, QuickCreateHorseComponent, QuickCreateAppointmentComponent],
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss'],
})
export class AppComponent {
  private readonly router = inject(Router);
  readonly showLayout = signal(true);
  readonly sidebarOpen = signal(false);

  constructor() {
    this.router.events
      .pipe(filter((e) => e instanceof NavigationEnd))
      .subscribe((e) => {
        const event = e as NavigationEnd;
        this.showLayout.set(!event.urlAfterRedirects.startsWith('/auth'));
        this.sidebarOpen.set(false);
      });
  }

  toggleSidebar(): void {
    this.sidebarOpen.update((v) => !v);
  }
}
