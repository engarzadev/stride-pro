import { Component, inject, computed, signal, HostListener } from '@angular/core';
import { RouterLink, RouterLinkActive, Router } from '@angular/router';
import { MatIconModule } from '@angular/material/icon';
import { toSignal } from '@angular/core/rxjs-interop';
import { AuthService } from '../../../core/services/auth.service';

interface BottomNavItem {
  path: string;
  label: string;
  icon: string;
  isCenter?: boolean;
}

interface QuickAction {
  path: string;
  label: string;
  icon: string;
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
  private readonly router = inject(Router);
  private readonly currentUser = toSignal(this.authService.currentUser$);

  readonly addMenuOpen = signal(false);

  readonly quickActions: QuickAction[] = [
    { path: '/horses/new', label: 'Add Horse', icon: 'pets' },
    { path: '/care-log', label: 'Care Log Entry', icon: 'monitor_heart' },
    { path: '/reminders', label: 'Add Reminder', icon: 'notifications_active' },
  ];

  private readonly ownerItems: BottomNavItem[] = [
    { path: '/dashboard', label: 'Home', icon: 'home' },
    { path: '/horses', label: 'Horses', icon: 'pets' },
    { path: '', label: 'Add', icon: 'add', isCenter: true },
    { path: '/reminders', label: 'Reminders', icon: 'notifications' },
    { path: '/settings', label: 'Profile', icon: 'person' },
  ];

  private readonly professionalItems: BottomNavItem[] = [
    { path: '/dashboard', label: 'Home', icon: 'home' },
    { path: '/clients', label: 'Clients', icon: 'people' },
    { path: '', label: 'Add', icon: 'add', isCenter: true },
    { path: '/horses', label: 'Horses', icon: 'pets' },
    { path: '/settings', label: 'Profile', icon: 'person' },
  ];

  readonly navItems = computed(() => {
    const user = this.currentUser();
    return user?.role === 'owner' ? this.ownerItems : this.professionalItems;
  });

  toggleAddMenu(event: Event): void {
    event.stopPropagation();
    this.addMenuOpen.update((v) => !v);
  }

  navigateTo(path: string): void {
    this.addMenuOpen.set(false);
    this.router.navigate([path]);
  }

  @HostListener('document:click')
  closeMenu(): void {
    this.addMenuOpen.set(false);
  }
}
