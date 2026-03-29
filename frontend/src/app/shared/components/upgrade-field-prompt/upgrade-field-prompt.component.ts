import { Component, input } from '@angular/core';
import { RouterLink } from '@angular/router';
import { MatIconModule } from '@angular/material/icon';

@Component({
  selector: 'app-upgrade-field-prompt',
  standalone: true,
  imports: [RouterLink, MatIconModule],
  templateUrl: './upgrade-field-prompt.component.html',
  styleUrls: ['./upgrade-field-prompt.component.scss'],
})
export class UpgradeFieldPromptComponent {
  readonly label = input.required<string>();
}
