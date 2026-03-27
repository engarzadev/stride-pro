import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'currencyFormat',
  standalone: true,
})
export class CurrencyFormatPipe implements PipeTransform {
  transform(value: unknown, currency: string = 'USD'): string {
    if (value == null) return '';
    const num = Number(value);
    if (isNaN(num)) return String(value);
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency,
      minimumFractionDigits: 2,
    }).format(num);
  }
}
