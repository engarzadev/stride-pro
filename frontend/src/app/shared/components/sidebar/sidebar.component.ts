import { Component, Input, Output, EventEmitter } from '@angular/core';
import { RouterLink, RouterLinkActive } from '@angular/router';
import { MatIconModule } from '@angular/material/icon';

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

  readonly navItems = [
    { path: '/dashboard', label: 'Dashboard', icon: 'home' },
    { path: '/clients', label: 'Clients', icon: 'people' },
    { path: '/horses', label: 'Horses', icon: 'pets' },
    { path: '/barns', label: 'Barns', icon: 'warehouse' },
    { path: '/appointments', label: 'Appointments', icon: 'event' },
    { path: '/sessions', label: 'Sessions', icon: 'medical_services' },
    { path: '/invoices', label: 'Invoices', icon: 'receipt' },
    { path: '/billing', label: 'Billing', icon: 'payments' },
  ];
}
