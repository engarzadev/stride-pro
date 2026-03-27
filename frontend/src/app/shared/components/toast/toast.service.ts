import { Injectable, signal } from '@angular/core';

export interface Toast {
  id: number;
  message: string;
  type: 'success' | 'error' | 'warning' | 'info';
}

@Injectable({ providedIn: 'root' })
export class ToastService {
  private nextId = 0;
  readonly toasts = signal<Toast[]>([]);

  success(message: string): void {
    this.show(message, 'success');
  }

  error(message: string): void {
    this.show(message, 'error');
  }

  warning(message: string): void {
    this.show(message, 'warning');
  }

  info(message: string): void {
    this.show(message, 'info');
  }

  remove(id: number): void {
    this.toasts.update((toasts) => toasts.filter((t) => t.id !== id));
  }

  private show(message: string, type: Toast['type']): void {
    const id = this.nextId++;
    this.toasts.update((toasts) => [...toasts, { id, message, type }]);
    setTimeout(() => this.remove(id), 5000);
  }
}
