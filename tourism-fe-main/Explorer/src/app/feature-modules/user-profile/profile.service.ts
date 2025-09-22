import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs';
import { User } from './profile/model/profile.model';
import { TokenStorage } from 'src/app/infrastructure/auth/jwt/token.service';

@Injectable({
  providedIn: 'root',
})
export class ProfileService {
  private apiURL = 'http://localhost:8083/api/v1/profile';

  constructor(private http: HttpClient, private tokenStorage: TokenStorage) {}

  getProfile(): Observable<User> {
    const token = this.tokenStorage.getAccessToken();

    if (!token) {
      console.error('No token found!');
      throw new Error('No token found');
    }

    const headers = new HttpHeaders().set('Authorization', `Bearer ${token}`);
    return this.http.get<User>(this.apiURL, { headers });
  }

  updateProfile(profileData: User): Observable<User> {
    const token = this.tokenStorage.getAccessToken();

    if (!token) {
      console.error('No token found!');
      throw new Error('No token found');
    }

    const headers = new HttpHeaders().set('Authorization', `Bearer ${token}`);
    return this.http.put<User>(this.apiURL, profileData, { headers });
  }
}