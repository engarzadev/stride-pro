import { Component, Input, Output, EventEmitter } from '@angular/core';
import { RouterLink, RouterLinkActive } from '@angular/router';

@Component({
  selector: 'app-sidebar',
  standalone: true,
  imports: [RouterLink, RouterLinkActive],
  templateUrl: './sidebar.component.html',
  styleUrls: ['./sidebar.component.scss'],
})
export class SidebarComponent {
  @Input() collapsed = false;
  @Output() closeSidebar = new EventEmitter<void>();

  readonly navItems = [
    { path: '/dashboard', label: 'Dashboard', icon: '\u2302' },
    { path: '/clients', label: 'Clients', icon: '\u263A' },
    { path: '/horses', label: 'Horses', icon: '\u2658' },
    { path: '/barns', label: 'Barns', icon: '\u2616' },
    { path: '/appointments', label: 'Appointments', icon: '\u2637' },
    { path: '/sessions', label: 'Sessions', icon: '\u2695' },
    { path: '/invoices', label: 'Invoices', icon: '\u2709' },
  ];
}
