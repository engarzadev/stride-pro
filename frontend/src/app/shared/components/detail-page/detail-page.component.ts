import { Component, EventEmitter, Input, Output } from '@angular/core';
import { RouterLink } from '@angular/router';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';

@Component({
  selector: 'app-detail-page',
  standalone: true,
  imports: [RouterLink, MatButtonModule, MatIconModule],
  templateUrl: './detail-page.component.html',
  styleUrls: ['./detail-page.component.scss'],
})
export class DetailPageComponent {
  @Input() title = '';
  @Input() backRoute = '/';
  @Input() backLabel = 'Back';
  @Input() editRoute: string[] = [];
  @Input() extraActionLabel = '';
  @Input() extraActionIcon = '';
  @Input() extraActionDisabled = false;
  @Output() deleteClick = new EventEmitter<void>();
  @Output() extraActionClick = new EventEmitter<void>();
}
