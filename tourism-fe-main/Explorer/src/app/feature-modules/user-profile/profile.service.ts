import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs';
import { UpdateUserProfilePayload, User } from './profile/model/profile.model';
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

  updateProfile(payload: UpdateUserProfilePayload): Observable<User> {
    return this.http.put<User>(this.apiURL, payload, { headers: this.createAuthHeaders() });
  }

  //Kreira i vraća HttpHeaders sa JWT tokenom
  private createAuthHeaders(): HttpHeaders {    
    const token = this.tokenStorage.getAccessToken();

    if (!token) {
      console.error('No token found!');
      //backend će vratiti 401 Unauthorized gresku
      return new HttpHeaders();
    }
    
    // Vracamo hedere spremne za slanje
    return new HttpHeaders({
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`
    });
  }
}