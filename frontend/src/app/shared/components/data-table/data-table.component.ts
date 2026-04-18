import { Component, Input, Output, EventEmitter, signal, inject } from '@angular/core';
import { toSignal } from '@angular/core/rxjs-interop';
import { BreakpointObserver } from '@angular/cdk/layout';
import { map } from 'rxjs/operators';
import { MatTableModule } from '@angular/material/table';
import { MatSortModule, Sort } from '@angular/material/sort';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatChipsModule } from '@angular/material/chips';
import { MatTooltipModule } from '@angular/material/tooltip';
import { MatMenuModule } from '@angular/material/menu';
import { MatPaginatorModule, PageEvent } from '@angular/material/paginator';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';
import { CurrencyFormatPipe } from '../../pipes/currency-format.pipe';
import { DateFormatPipe } from '../../pipes/date-format.pipe';

export interface TableColumn {
  key: string;
  label: string;
  sortable?: boolean;
  type?: 'text' | 'date' | 'currency' | 'badge';
  badgeMap?: Record<string, string>;
  capitalize?: boolean;
}

export interface TableAction {
  label: string;
  icon?: string;
  class?: string;
  action: string;
}

export interface MobileCardConfig {
  /** Column key to display as the card's primary heading */
  titleKey: string;
  /** Optional column key to display as a secondary line below the title */
  subtitleKey?: string;
}

export interface FilterConfig {
  key: string;
  label: string;
  options: { value: string; label: string }[];
}

@Component({
  selector: 'app-data-table',
  standalone: true,
  imports: [MatTableModule, MatSortModule, MatButtonModule, MatIconModule, MatChipsModule, MatTooltipModule, MatMenuModule, MatPaginatorModule, MatFormFieldModule, MatInputModule, MatSelectModule, CurrencyFormatPipe, DateFormatPipe],
  templateUrl: './data-table.component.html',
  styleUrls: ['./data-table.component.scss'],
})
export class DataTableComponent {
  @Input() columns: TableColumn[] = [];
  @Input() data: Record<string, unknown>[] = [];
  @Input() actions: TableAction[] = [];
  @Input() filterConfig: FilterConfig[] = [];
  @Input() clickable = true;
  @Input() mobileCard?: MobileCardConfig;
  @Output() rowClick = new EventEmitter<Record<string, unknown>>();
  @Output() actionClick = new EventEmitter<{ action: string; row: Record<string, unknown> }>();

  private breakpointObserver = inject(BreakpointObserver);

  isMobile = toSignal(
    this.breakpointObserver.observe('(max-width: 640px)').pipe(map(r => r.matches)),
    { initialValue: false }
  );

  sortKey = signal('');
  sortDir = signal<'asc' | 'desc'>('asc');
  currentPage = signal(0);
  pageSize = signal(10);
  searchQuery = signal('');
  activeFilters = signal<Record<string, string>>({});

  get hasActiveSearch(): boolean {
    return !!this.searchQuery() || Object.values(this.activeFilters()).some(v => !!v);
  }

  get filteredCount(): number {
    return this.getFilteredData().length;
  }

  onSearchChange(query: string): void {
    this.searchQuery.set(query);
    this.currentPage.set(0);
  }

  setFilter(key: string, value: string): void {
    this.activeFilters.update(f => ({ ...f, [key]: value }));
    this.currentPage.set(0);
  }

  clearAllFilters(): void {
    this.searchQuery.set('');
    this.activeFilters.set({});
    this.currentPage.set(0);
  }

  getFilteredData(): Record<string, unknown>[] {
    let result = this.data;
    const query = this.searchQuery().toLowerCase().trim();
    if (query) {
      result = result.filter(row =>
        this.columns.some(col => {
          const val = this.getNestedValue(row, col.key);
          return val != null && String(val).toLowerCase().includes(query);
        })
      );
    }
    const filters = this.activeFilters();
    for (const [key, value] of Object.entries(filters)) {
      if (!value) continue;
      result = result.filter(row =>
        String(this.getNestedValue(row, key)).toLowerCase() === value.toLowerCase()
      );
    }
    return result;
  }

  get displayedColumns(): string[] {
    const cols = this.columns.map((c) => c.key);
    if (this.actions.length > 0) cols.push('_actions');
    return cols;
  }

  get mobileBodyColumns(): TableColumn[] {
    if (!this.mobileCard) return this.columns;
    const exclude = new Set([this.mobileCard.titleKey, this.mobileCard.subtitleKey].filter(Boolean) as string[]);
    return this.columns.filter(c => !exclude.has(c.key));
  }

  get mobileTitleColumn(): TableColumn | undefined {
    return this.columns.find(c => c.key === this.mobileCard?.titleKey);
  }

  get mobileSubtitleColumn(): TableColumn | undefined {
    return this.columns.find(c => c.key === this.mobileCard?.subtitleKey);
  }

  onSort(sort: Sort): void {
    this.sortKey.set(sort.active);
    this.sortDir.set(sort.direction as 'asc' | 'desc' || 'asc');
    this.currentPage.set(0);
  }

  onPage(event: PageEvent): void {
    this.currentPage.set(event.pageIndex);
    this.pageSize.set(event.pageSize);
  }

  getPagedData(): Record<string, unknown>[] {
    const start = this.currentPage() * this.pageSize();
    return this.getSortedData().slice(start, start + this.pageSize());
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
    const filtered = this.getFilteredData();
    if (!key) return filtered;
    const dir = this.sortDir() === 'asc' ? 1 : -1;
    return [...filtered].sort((a, b) => {
      const aVal = this.getNestedValue(a, key);
      const bVal = this.getNestedValue(b, key);
      if (aVal == null) return 1;
      if (bVal == null) return -1;
      if (aVal < bVal) return -1 * dir;
      if (aVal > bVal) return 1 * dir;
      return 0;
    });
  }

  getChipColor(badgeClass: string): string {
    if (badgeClass === 'primary') return 'primary';
    if (badgeClass === 'danger') return 'warn';
    return '';
  }
}
