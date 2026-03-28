import { Component, forwardRef, signal, computed } from '@angular/core';
import { ControlValueAccessor, NG_VALUE_ACCESSOR, ReactiveFormsModule } from '@angular/forms';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatAutocompleteModule } from '@angular/material/autocomplete';

const HORSE_BREEDS = [
  'American Quarter Horse',
  'American Saddlebred',
  'Andalusian',
  'Appaloosa',
  'Arabian',
  'Belgian Draft',
  'Clydesdale',
  'Connemara Pony',
  'Dutch Warmblood (KWPN)',
  'Fjord',
  'Friesian',
  'Haflinger',
  'Hanoverian',
  'Icelandic Horse',
  'Lipizzaner',
  'Lusitano',
  'Miniature Horse',
  'Missouri Fox Trotter',
  'Morgan',
  'Mustang',
  'Oldenburg',
  'Paint Horse',
  'Paso Fino',
  'Percheron',
  'Rocky Mountain Horse',
  'Shetland Pony',
  'Shire',
  'Spotted Saddle Horse',
  'Standardbred',
  'Tennessee Walking Horse',
  'Thoroughbred',
  'Trakehner',
  'Warmblood',
  'Welsh Pony',
];

@Component({
  selector: 'app-breed-autocomplete',
  standalone: true,
  imports: [ReactiveFormsModule, MatFormFieldModule, MatInputModule, MatAutocompleteModule],
  templateUrl: './breed-autocomplete.component.html',
  styleUrls: ['./breed-autocomplete.component.scss'],
  providers: [
    {
      provide: NG_VALUE_ACCESSOR,
      useExisting: forwardRef(() => BreedAutocompleteComponent),
      multi: true,
    },
  ],
})
export class BreedAutocompleteComponent implements ControlValueAccessor {
  readonly inputValue = signal('');
  readonly isDisabled = signal(false);

  readonly filteredBreeds = computed(() => {
    const query = this.inputValue().toLowerCase().trim();
    if (!query) return HORSE_BREEDS;
    return HORSE_BREEDS.filter((b) => b.toLowerCase().includes(query));
  });

  private onChange: (value: string) => void = () => {};
  private onTouched: () => void = () => {};

  writeValue(value: string): void {
    this.inputValue.set(value ?? '');
  }

  registerOnChange(fn: (value: string) => void): void {
    this.onChange = fn;
  }

  registerOnTouched(fn: () => void): void {
    this.onTouched = fn;
  }

  setDisabledState(isDisabled: boolean): void {
    this.isDisabled.set(isDisabled);
  }

  onInput(value: string): void {
    this.inputValue.set(value);
    this.onChange(value);
  }

  onBlur(): void {
    this.onTouched();
  }

  selectBreed(breed: string): void {
    this.inputValue.set(breed);
    this.onChange(breed);
  }
}
