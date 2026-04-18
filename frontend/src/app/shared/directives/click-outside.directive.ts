import { Directive, ElementRef, Output, EventEmitter, inject, OnInit, OnDestroy } from '@angular/core';

@Directive({
  selector: '[appClickOutside]',
  standalone: true,
})
export class ClickOutsideDirective implements OnInit, OnDestroy {
  private readonly el = inject(ElementRef);
  @Output() clickOutside = new EventEmitter<void>();

  private handler = (event: Event): void => {
    if (!this.el.nativeElement.contains(event.target)) {
      this.clickOutside.emit();
    }
  };

  ngOnInit(): void {
    setTimeout(() => {
      document.addEventListener('click', this.handler);
    });
  }

  ngOnDestroy(): void {
    document.removeEventListener('click', this.handler);
  }
}
