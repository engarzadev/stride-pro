import { Component, Injectable, signal } from '@angular/core';

export interface ConfirmDialogData {
  title: string;
  message: string;
  confirmText?: string;
  cancelText?: string;
  confirmClass?: string;
}

@Injectable({ providedIn: 'root' })
export class ConfirmDialogService {
  readonly visible = signal(false);
  readonly data = signal<ConfirmDialogData>({ title: '', message: '' });

  private resolveFn?: (result: boolean) => void;

  confirm(data: ConfirmDialogData): Promise<boolean> {
    this.data.set(data);
    this.visible.set(true);
    return new Promise<boolean>((resolve) => {
      this.resolveFn = resolve;
    });
  }

  accept(): void {
    this.visible.set(false);
    this.resolveFn?.(true);
  }

  cancel(): void {
    this.visible.set(false);
    this.resolveFn?.(false);
  }
}

@Component({
  selector: 'app-confirm-dialog',
  standalone: true,
  templateUrl: './confirm-dialog.component.html',
  styleUrls: ['./confirm-dialog.component.scss'],
})
export class ConfirmDialogComponent {
  constructor(readonly dialogService: ConfirmDialogService) {}
}
