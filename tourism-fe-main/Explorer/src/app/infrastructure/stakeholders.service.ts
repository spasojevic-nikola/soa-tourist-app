import { Injectable } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from 'src/env/environment';
import { UpdateUserProfilePayload, User } from '../feature-modules/user-profile/profile/model/profile.model';

@Injectable({
  providedIn: 'root'
})
export class StakeholdersService {
  
  constructor(private http: HttpClient) { }


  searchUsers(username: string): Observable<User[]> {
    const params = new HttpParams().set('username', username);
    return this.http.get<User[]>(`${environment.stakeholdersApiHost}/users/search`, { params });
  }

 /*
  getProfile(): Observable<User> {
    return this.http.get<User>(`${environment.stakeholdersApiHost}profile`);
  }*/

/*
  updateProfile(payload: UpdateUserProfilePayload): Observable<User> {
    return this.http.put<User>(`${environment.stakeholdersApiHost}profile`, payload);
  }*/


  getUserById(id: number): Observable<User> {
    return this.http.get<User>(`${environment.stakeholdersApiHost}/users/${id}`);
  }
  
}