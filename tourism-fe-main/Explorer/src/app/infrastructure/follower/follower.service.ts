import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { environment } from 'src/env/environment';

@Injectable({
  providedIn: 'root'
})
export class FollowerService {

  private apiBaseUrl = environment.followerApiHost;

  constructor(private http: HttpClient) { }

  checkIfFollowing(userId: number): Observable<{ follows: boolean }> {
    // Hederi se dodaju automatski preko interceptora
    return this.http.get<{ follows: boolean }>(`${this.apiBaseUrl}/check-follow/${userId}`);
  }

  follow(userId: number): Observable<any> {
    return this.http.post<any>(`${this.apiBaseUrl}/follow/${userId}`, {});
  }

  unfollow(userId: number): Observable<any> {
    return this.http.delete<any>(`${this.apiBaseUrl}/unfollow/${userId}`);
  }

}
