import { HttpInterceptorFn } from '@angular/common/http';
import { inject } from '@angular/core';
import { Router } from '@angular/router';
import { catchError, throwError } from 'rxjs';
import { ToastService } from '../../shared/components/toast/toast.service';

export const errorInterceptor: HttpInterceptorFn = (req, next) => {
  const router = inject(Router);
  const toast = inject(ToastService);

  return next(req).pipe(
    catchError((error) => {
      let message = 'An unexpected error occurred';

      if (error.status === 0) {
        message = 'Unable to connect to server';
      } else if (error.status === 401) {
        const isAuthRequest = req.url.includes('/auth/login') || req.url.includes('/auth/register');
        if (isAuthRequest) {
          message = error.error?.error?.message || 'Invalid email or password';
        } else {
          message = 'Your session has expired. Please log in again.';
          localStorage.removeItem('stride_pro_user');
          router.navigate(['/auth/login']);
        }
      } else if (error.status === 403) {
        const errMsg = error.error?.error?.message || error.error?.error || '';
        if (typeof errMsg === 'string' && errMsg.toLowerCase().includes('require')) {
          return throwError(() => error);
        }
        message = 'You do not have permission to perform this action';
      } else if (error.status === 404) {
        message = 'The requested resource was not found';
      } else if (error.status === 409) {
        message = error.error?.error?.message || 'Scheduling conflict detected';
      } else if (error.status === 422) {
        const err = error.error?.error;
        const fields = err?.fields;
        if (fields && typeof fields === 'object') {
          message = Object.values(fields).join(', ');
        } else {
          message = err?.message || (typeof err === 'string' ? err : 'Validation error');
        }
      } else if (error.status >= 500) {
        message = 'Server error. Please try again later.';
      } else if (error.error?.error) {
        const err = error.error.error;
        message = typeof err === 'string' ? err : err?.message || message;
      }

      message = message.charAt(0).toUpperCase() + message.slice(1);
      toast.error(message);
      return throwError(() => error);
    })
  );
};
