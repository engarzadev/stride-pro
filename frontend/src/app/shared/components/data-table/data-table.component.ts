import { Component, Input, Output, EventEmitter, signal } from '@angular/core';
import { CurrencyFormatPipe } from '../../pipes/currency-format.pipe';
import { DateFormatPipe } from '../../pipes/date-format.pipe';

export interface TableColumn {
  key: string;
  label: string;
  sortable?: boolean;
  type?: 'text' | 'date' | 'currency' | 'badge';
  badgeMap?: Record<string, string>;
}

export interface TableAction {
  label: string;
  icon?: string;
  class?: string;
  action: string;
}

@Component({
  selector: 'app-data-table',
  standalone: true,
  imports: [CurrencyFormatPipe, DateFormatPipe],
  templateUrl: './data-table.component.html',
  styleUrls: ['./data-table.component.scss'],
})
export class DataTableComponent {
  @Input() columns: TableColumn[] = [];
  @Input() data: Record<string, unknown>[] = [];
  @Input() actions: TableAction[] = [];
  @Input() clickable = true;
  @Output() rowClick = new EventEmitter<Record<string, unknown>>();
  @Output() actionClick = new EventEmitter<{ action: string; row: Record<string, unknown> }>();

  sortKey = signal('');
  sortDir = signal<'asc' | 'desc'>('asc');

  onSort(column: TableColumn): void {
    if (!column.sortable) return;
    if (this.sortKey() === column.key) {
      this.sortDir.update((d) => (d === 'asc' ? 'desc' : 'asc'));
    } else {
      this.sortKey.set(column.key);
      this.sortDir.set('asc');
    }
  }

  getNestedValue(obj: Record<string, unknown>, path: string): unknown {
    return path.split('.').reduce((o: unknown, k: string) => {
      if (o && typeof o === 'object') {
        return (o as Record<string, unknown>)[k];
      }
      return undefined;
    }, obj);
  }

  getSortedData(): Record<string, unknown>[] {
    const key = this.sortKey();
    if (!key) return this.data;
    const dir = this.sortDir() === 'asc' ? 1 : -1;
    return [...this.data].sort((a, b) => {
      const aVal = this.getNestedValue(a, key);
      const bVal = this.getNestedValue(b, key);
      if (aVal == null) return 1;
      if (bVal == null) return -1;
      if (aVal < bVal) return -1 * dir;
      if (aVal > bVal) return 1 * dir;
      return 0;
    });
  }
}
