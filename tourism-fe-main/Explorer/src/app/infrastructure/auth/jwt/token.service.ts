import { Injectable } from '@angular/core';
import { ACCESS_TOKEN , USER } from '../../../shared/constants';

@Injectable({
    providedIn: 'root',
  })
  export class TokenStorage {
    constructor() {}
  
    saveAccessToken(token: string): void {
      localStorage.removeItem(ACCESS_TOKEN);
      localStorage.setItem(ACCESS_TOKEN, token);
    }
  
    getAccessToken() {
      return localStorage.getItem(ACCESS_TOKEN);
    }
  
    clear() {
      localStorage.removeItem(ACCESS_TOKEN);
      localStorage.removeItem(USER);
    }
    
    getUserId(): number | null {
    const token = this.getAccessToken();
    if (!token) return null;

    try {
      const payload = JSON.parse(atob(token.split('.')[1]));
      // backend mora da šalje userId ili sub u tokenu
      return payload.userId ?? payload.sub ?? null;
    } catch (err) {
      console.error('Neispravan token', err);
      return null;
    }
  }
  }