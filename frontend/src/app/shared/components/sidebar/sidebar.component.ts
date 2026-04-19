import { Component, Input, Output, EventEmitter, inject, computed } from '@angular/core';
import { RouterLink, RouterLinkActive } from '@angular/router';
import { MatIconModule } from '@angular/material/icon';
import { toSignal } from '@angular/core/rxjs-interop';
import { AuthService } from '../../../core/services/auth.service';

interface NavItem {
  path: string;
  label: string;
  icon: string;
  roles: string[];
}

@Component({
  selector: 'app-sidebar',
  standalone: true,
  imports: [RouterLink, RouterLinkActive, MatIconModule],
  templateUrl: './sidebar.component.html',
  styleUrls: ['./sidebar.component.scss'],
})
export class SidebarComponent {
  @Input() collapsed = false;
  @Output() closeSidebar = new EventEmitter<void>();

  private readonly authService = inject(AuthService);
  private readonly currentUser = toSignal(this.authService.currentUser$);

  private readonly allNavItems: NavItem[] = [
    { path: '/dashboard', label: 'Dashboard', icon: 'home', roles: ['owner', 'professional'] },
    { path: '/clients', label: 'Clients', icon: 'people', roles: ['professional'] },
    { path: '/horses', label: 'Horses', icon: 'pets', roles: ['owner', 'professional'] },
    { path: '/care-log', label: 'Care Log', icon: 'monitor_heart', roles: ['owner'] },
    { path: '/reminders', label: 'Reminders', icon: 'notifications_active', roles: ['owner'] },
    { path: '/barns', label: 'Barns', icon: 'warehouse', roles: ['professional'] },
    { path: '/appointments', label: 'Appointments', icon: 'event', roles: ['professional'] },
    { path: '/sessions', label: 'Sessions', icon: 'medical_services', roles: ['professional'] },
    { path: '/invoices', label: 'Invoices', icon: 'receipt', roles: ['professional'] },
    { path: '/billing', label: 'Billing', icon: 'payments', roles: ['professional'] },
    { path: '/settings', label: 'Settings', icon: 'settings', roles: ['owner', 'professional'] },
  ];

  readonly navItems = computed(() => {
    const user = this.currentUser();
    // Default to showing all items if role is unknown (backward compat for "user" role)
    const role = user?.role === 'owner' ? 'owner' : 'professional';
    return this.allNavItems.filter(item => item.roles.includes(role));
  });
}
