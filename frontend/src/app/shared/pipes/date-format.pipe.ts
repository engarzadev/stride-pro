import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'dateFormat',
  standalone: true,
})
export class DateFormatPipe implements PipeTransform {
  transform(value: unknown, format: 'short' | 'long' | 'time' = 'short'): string {
    if (!value) return '';
    const date = new Date(value as string);
    if (isNaN(date.getTime())) return String(value);

    switch (format) {
      case 'long':
        return date.toLocaleDateString('en-US', {
          weekday: 'long',
          year: 'numeric',
          month: 'long',
          day: 'numeric',
        });
      case 'time':
        return date.toLocaleString('en-US', {
          month: 'short',
          day: 'numeric',
          year: 'numeric',
          hour: 'numeric',
          minute: '2-digit',
        });
      case 'short':
      default:
        return date.toLocaleDateString('en-US', {
          month: 'short',
          day: 'numeric',
          year: 'numeric',
        });
    }
  }
}
