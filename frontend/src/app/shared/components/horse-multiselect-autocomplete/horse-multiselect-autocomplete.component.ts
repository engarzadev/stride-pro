import { Component, ElementRef, EventEmitter, Input, OnChanges, Output, ViewChild, computed, signal } from '@angular/core';
import { MatAutocompleteModule, MatAutocompleteSelectedEvent } from '@angular/material/autocomplete';
import { MatChipsModule } from '@angular/material/chips';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { Horse } from '../../../core/models';

@Component({
  selector: 'app-horse-multiselect-autocomplete',
  standalone: true,
  imports: [MatFormFieldModule, MatInputModule, MatChipsModule, MatAutocompleteModule, MatIconModule],
  templateUrl: './horse-multiselect-autocomplete.component.html',
})
export class HorseMultiselectAutocompleteComponent implements OnChanges {
  @Input() horses: Horse[] = [];
  @Input() selectedHorses: Horse[] = [];
  @Input() disabled = false;
  @Output() selectedHorsesChange = new EventEmitter<Horse[]>();

  @ViewChild('searchInput') searchInput!: ElementRef<HTMLInputElement>;

  readonly _horses = signal<Horse[]>([]);
  readonly _selectedHorses = signal<Horse[]>([]);
  readonly searchQuery = signal('');

  readonly displayedSelectedHorses = computed(() => this._selectedHorses());

  readonly filteredHorses = computed(() => {
    const query = this.searchQuery().toLowerCase().trim();
    const ids = new Set(this._selectedHorses().map((h) => h.id));
    return this._horses()
      .filter((h) => !ids.has(h.id))
      .filter((h) => !query || h.name.toLowerCase().includes(query));
  });

  ngOnChanges(): void {
    this._horses.set(this.horses ?? []);
    this._selectedHorses.set(this.selectedHorses ?? []);
  }

  select(event: MatAutocompleteSelectedEvent): void {
    const horse: Horse = event.option.value;
    this._selectedHorses.update((list) => [...list, horse]);
    this.searchQuery.set('');
    this.searchInput.nativeElement.value = '';
    this.selectedHorsesChange.emit(this._selectedHorses());
  }

  remove(horseId: string): void {
    this._selectedHorses.update((list) => list.filter((h) => h.id !== horseId));
    this.selectedHorsesChange.emit(this._selectedHorses());
  }

  onInput(value: string): void {
    this.searchQuery.set(value);
  }
}
