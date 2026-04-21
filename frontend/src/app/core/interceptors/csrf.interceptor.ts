import { HttpInterceptorFn } from '@angular/common/http';

const CSRF_COOKIE = 'XSRF-TOKEN';
const CSRF_HEADER = 'X-XSRF-TOKEN';
const SAFE_METHODS = ['GET', 'HEAD', 'OPTIONS'];

// Angular's built-in withXsrfConfiguration skips absolute URLs, so we need a
// custom interceptor for cross-origin API requests.
export const csrfInterceptor: HttpInterceptorFn = (req, next) => {
  if (SAFE_METHODS.includes(req.method)) {
    return next(req);
  }

  const token = readCookie(CSRF_COOKIE);
  if (token) {
    return next(req.clone({ setHeaders: { [CSRF_HEADER]: token } }));
  }

  return next(req);
};

function readCookie(name: string): string | null {
  const match = document.cookie.match(new RegExp(`(?:^|;\\s*)${name}=([^;]*)`));
  return match ? decodeURIComponent(match[1]) : null;
}
