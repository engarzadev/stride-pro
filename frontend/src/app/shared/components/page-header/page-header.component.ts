import { Component, Input, Output, EventEmitter } from '@angular/core';

@Component({
  selector: 'app-page-header',
  standalone: true,
  templateUrl: './page-header.component.html',
  styleUrls: ['./page-header.component.scss'],
})
export class PageHeaderComponent {
  @Input() title = '';
  @Input() buttonText = '';
  @Output() buttonClick = new EventEmitter<void>();
}
