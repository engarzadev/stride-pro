import { Component, computed, HostListener, inject, signal } from '@angular/core';
import { toSignal } from '@angular/core/rxjs-interop';
import { MatIconModule } from '@angular/material/icon';
import { Router, RouterLink, RouterLinkActive } from '@angular/router';
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

  private readonly ownerActions: QuickAction[] = [
    { path: '/horses/new', label: 'Add Horse', icon: 'chess_knight' },
    { path: '/care-log?showForm=true', label: 'Care Log Entry', icon: 'monitor_heart' },
    { path: '/reminders?showForm=true', label: 'Add Reminder', icon: 'notifications_active' },
  ];

  private readonly professionalActions: QuickAction[] = [
    { path: '/clients/new', label: 'New Client', icon: 'person_add' },
    { path: '/appointments/new', label: 'New Appointment', icon: 'event' },
    { path: '/horses/new', label: 'Add Horse', icon: 'chess_knight' },
    { path: '/invoices/new', label: 'Add Invoice', icon: 'receipt_long' },
    { path: '/sessions/new', label: 'Add Session', icon: 'event' },
  ];

  readonly quickActions = computed(() => {
    const user = this.currentUser();
    return user?.role === 'owner' ? this.ownerActions : this.professionalActions;
  });

  private readonly ownerItems: BottomNavItem[] = [
    { path: '/dashboard', label: 'Home', icon: 'home' },
    { path: '/horses', label: 'Horses', icon: 'chess_knight' },
    { path: '', label: 'Add', icon: 'add', isCenter: true },
    { path: '/reminders', label: 'Reminders', icon: 'notifications' },
    { path: '/settings', label: 'Settings', icon: 'settings' },
  ];

  private readonly professionalItems: BottomNavItem[] = [
    { path: '/dashboard', label: 'Home', icon: 'home' },
    { path: '/invoices', label: 'Invoices', icon: 'receipt_long' },
    // { path: '/clients', label: 'Clients', icon: 'people' },
    // { path: '/horses', label: 'Horses', icon: 'chess_knight' },
    { path: '', label: 'Add', icon: 'add', isCenter: true },
    { path: '/sessions', label: 'Sessions', icon: 'event' },
    { path: '/settings', label: 'Settings', icon: 'settings' },
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
    this.router.navigateByUrl(path);
  }

  @HostListener('document:click')
  closeMenu(): void {
    this.addMenuOpen.set(false);
  }
}
