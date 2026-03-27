import { Component, inject, signal, Output, EventEmitter } from '@angular/core';
import { Router } from '@angular/router';
import { AsyncPipe } from '@angular/common';
import { AuthService } from '../../../core/services/auth.service';
import { ClickOutsideDirective } from '../../directives/click-outside.directive';

@Component({
  selector: 'app-header',
  standalone: true,
  imports: [AsyncPipe, ClickOutsideDirective],
  templateUrl: './header.component.html',
  styleUrls: ['./header.component.scss'],
})
export class HeaderComponent {
  private readonly authService = inject(AuthService);
  private readonly router = inject(Router);

  @Output() toggleSidebar = new EventEmitter<void>();

  readonly currentUser$ = this.authService.currentUser$;
  readonly menuOpen = signal(false);

  toggleMenu(): void {
    this.menuOpen.update((v) => !v);
  }

  closeMenu(): void {
    this.menuOpen.set(false);
  }

  logout(): void {
    this.authService.logout();
    this.closeMenu();
  }
}
