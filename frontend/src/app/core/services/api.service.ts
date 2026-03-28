import { Injectable, inject } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';
import { environment } from '../../../environments/environment';
import { keysToCamel, keysToSnake } from '../utils/camel-case';

@Injectable({ providedIn: 'root' })
export class ApiService {
  private readonly http = inject(HttpClient);
  private readonly baseUrl = environment.apiUrl;

  get<T>(path: string, params?: Record<string, string | number | boolean>): Observable<T> {
    let httpParams = new HttpParams();
    if (params) {
      Object.entries(params).forEach(([key, value]) => {
        if (value !== undefined && value !== null) {
          httpParams = httpParams.set(key, String(value));
        }
      });
    }
    return this.http.get<{ data: T }>(`${this.baseUrl}${path}`, { params: httpParams }).pipe(map((r) => keysToCamel<T>(r.data)));
  }

  post<T>(path: string, body: unknown): Observable<T> {
    return this.http.post<{ data: T }>(`${this.baseUrl}${path}`, keysToSnake(body)).pipe(map((r) => keysToCamel<T>(r.data)));
  }

  put<T>(path: string, body: unknown): Observable<T> {
    return this.http.put<{ data: T }>(`${this.baseUrl}${path}`, keysToSnake(body)).pipe(map((r) => keysToCamel<T>(r.data)));
  }

  delete<T>(path: string): Observable<T> {
    return this.http.delete<{ data: T }>(`${this.baseUrl}${path}`).pipe(map((r) => keysToCamel<T>(r.data)));
  }
}
