import { Component, inject, computed } from '@angular/core';
import { RouterLink, RouterLinkActive } from '@angular/router';
import { MatIconModule } from '@angular/material/icon';
import { toSignal } from '@angular/core/rxjs-interop';
import { AuthService } from '../../../core/services/auth.service';

interface BottomNavItem {
  path: string;
  label: string;
  icon: string;
  isCenter?: boolean;
}

@Component({
  selector: 'app-bottom-nav',
  standalone: true,
  imports: [RouterLink, RouterLinkActive, MatIconModule],
  templateUrl: './bottom-nav.component.html',
  styleUrls: ['./bottom-nav.component.scss'],
})
export class BottomNavComponent {
  private readonly authService = inject(AuthService);
  private readonly currentUser = toSignal(this.authService.currentUser$);

  private readonly ownerItems: BottomNavItem[] = [
    { path: '/dashboard', label: 'Home', icon: 'home' },
    { path: '/horses', label: 'Horses', icon: 'pets' },
    { path: '/care-log', label: 'Add', icon: 'add', isCenter: true },
    { path: '/reminders', label: 'Reminders', icon: 'notifications' },
    { path: '/settings', label: 'Profile', icon: 'person' },
  ];

  private readonly professionalItems: BottomNavItem[] = [
    { path: '/dashboard', label: 'Home', icon: 'home' },
    { path: '/clients', label: 'Clients', icon: 'people' },
    { path: '/appointments', label: 'Add', icon: 'add', isCenter: true },
    { path: '/horses', label: 'Horses', icon: 'pets' },
    { path: '/settings', label: 'Profile', icon: 'person' },
  ];

  readonly navItems = computed(() => {
    const user = this.currentUser();
    return user?.role === 'owner' ? this.ownerItems : this.professionalItems;
  });
}
