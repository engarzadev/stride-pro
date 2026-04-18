import { Component, EventEmitter, Input, Output } from '@angular/core';
import { RouterLink } from '@angular/router';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { MatIconModule } from '@angular/material/icon';

@Component({
  selector: 'app-form-page',
  standalone: true,
  imports: [RouterLink, MatButtonModule, MatCardModule, MatIconModule],
  templateUrl: './form-page.component.html',
  styleUrls: ['./form-page.component.scss'],
})
export class FormPageComponent {
  @Input() title = '';
  @Input() backRoute: string | string[] = '/';
  @Input() backLabel = 'Back';
  @Input() saving = false;
  @Input() formInvalid = false;
  @Input() isEdit = false;
  @Input() createIcon = 'add';
  @Output() save = new EventEmitter<void>();
}
