import { Component, Input, Output, EventEmitter } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';

@Component({
  selector: 'app-page-header',
  standalone: true,
  imports: [MatButtonModule],
  templateUrl: './page-header.component.html',
  styleUrls: ['./page-header.component.scss'],
})
export class PageHeaderComponent {
  @Input() title = '';
  @Input() buttonText = '';
  @Input() secondButtonText = '';
  @Output() buttonClick = new EventEmitter<void>();
  @Output() secondButtonClick = new EventEmitter<void>();
}
