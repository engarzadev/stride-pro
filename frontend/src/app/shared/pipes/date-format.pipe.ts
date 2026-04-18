import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'dateFormat',
  standalone: true,
})
export class DateFormatPipe implements PipeTransform {
  transform(value: unknown, format: 'short' | 'long' | 'time' | 'numeric' = 'short'): string {
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
      case 'numeric': {
        // Parse YYYY-MM-DD directly to avoid UTC-to-local timezone shifts
        const raw = String(value).substring(0, 10);
        const parts = raw.split('-');
        if (parts.length === 3) return `${parts[1]}/${parts[2]}/${parts[0]}`;
        return raw;
      }
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
