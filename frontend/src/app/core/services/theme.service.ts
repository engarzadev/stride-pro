import { Injectable, signal } from '@angular/core';

@Injectable({ providedIn: 'root' })
export class ThemeService {
  private readonly isDarkSignal = signal(false);
  readonly isDark = this.isDarkSignal.asReadonly();

  constructor() {
    this.init();
  }

  toggle(): void {
    const next = !this.isDarkSignal();
    this.isDarkSignal.set(next);
    localStorage.setItem('theme', next ? 'dark' : 'light');
    this.applyTheme(next);
  }

  private init(): void {
    const saved = localStorage.getItem('theme');
    const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
    const dark = saved ? saved === 'dark' : prefersDark;
    this.isDarkSignal.set(dark);
    this.applyTheme(dark);
  }

  private applyTheme(dark: boolean): void {
    document.documentElement.classList.toggle('dark', dark);
  }
}
