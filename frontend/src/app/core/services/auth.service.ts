import { Injectable, inject } from '@angular/core';
import { Router } from '@angular/router';
import { BehaviorSubject, Observable, tap } from 'rxjs';
import { AuthResponse, LoginRequest, RegisterRequest, User } from '../models';
import { ApiService } from './api.service';
import { keysToCamel } from '../utils/camel-case';

interface UpdateProfileRequest {
  firstName: string;
  lastName: string;
  email: string;
}

@Injectable({ providedIn: 'root' })
export class AuthService {
  private readonly api = inject(ApiService);
  private readonly router = inject(Router);

  private readonly USER_KEY = 'stride_pro_user';

  private readonly _isAuthenticated$ = new BehaviorSubject<boolean>(
    this.hasStoredUser(),
  );
  private readonly _currentUser$ = new BehaviorSubject<User | null>(
    this.getStoredUser(),
  );

  readonly isAuthenticated$ = this._isAuthenticated$.asObservable();
  readonly currentUser$ = this._currentUser$.asObservable();

  login(email: string, password: string): Observable<AuthResponse> {
    const body: LoginRequest = { email, password };
    return this.api
      .post<AuthResponse>('/auth/login', body)
      .pipe(tap((response) => this.handleAuthResponse(response)));
  }

  register(
    email: string,
    password: string,
    firstName: string,
    lastName: string,
    accountType: string,
  ): Observable<AuthResponse> {
    const body: RegisterRequest = {
      email,
      password,
      first_name: firstName,
      last_name: lastName,
      account_type: accountType,
    };
    return this.api
      .post<AuthResponse>('/auth/register', body)
      .pipe(tap((response) => this.handleAuthResponse(response)));
  }

  updateProfile(firstName: string, lastName: string, email: string): Observable<User> {
    const body: UpdateProfileRequest = { firstName, lastName, email };
    return this.api.put<User>('/auth/profile', body).pipe(
      tap((user) => {
        localStorage.setItem(this.USER_KEY, JSON.stringify(user));
        this._currentUser$.next(user);
      }),
    );
  }

  changePassword(currentPassword: string, newPassword: string): Observable<void> {
    return this.api.post<void>('/auth/change-password', { currentPassword, newPassword });
  }

  logout(): void {
    // Tell the server to revoke the token and clear the HttpOnly cookies
    this.api.post('/auth/logout', {}).subscribe({
      complete: () => this.clearSession(),
      error: () => this.clearSession(),
    });
  }

  isAuthenticated(): boolean {
    return this.hasStoredUser();
  }

  getStoredUser(): User | null {
    const userJson = localStorage.getItem(this.USER_KEY);
    if (userJson) {
      try {
        return keysToCamel<User>(JSON.parse(userJson));
      } catch {
        return null;
      }
    }
    return null;
  }

  private handleAuthResponse(response: AuthResponse): void {
    // Tokens are set as HttpOnly cookies by the server — never stored in JS
    localStorage.setItem(this.USER_KEY, JSON.stringify(response.user));
    this._isAuthenticated$.next(true);
    this._currentUser$.next(response.user);
  }

  private clearSession(): void {
    localStorage.removeItem(this.USER_KEY);
    this._isAuthenticated$.next(false);
    this._currentUser$.next(null);
    this.router.navigate(['/auth/login']);
  }

  private hasStoredUser(): boolean {
    return !!localStorage.getItem(this.USER_KEY);
  }
}
