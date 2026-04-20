import { HttpInterceptorFn } from '@angular/common/http';

// Tokens are stored as HttpOnly cookies set by the server and are sent
// automatically with every same-site request. withCredentials: true ensures
// the browser includes those cookies on all API calls.
export const authInterceptor: HttpInterceptorFn = (req, next) => {
  return next(req.clone({ withCredentials: true }));
};
