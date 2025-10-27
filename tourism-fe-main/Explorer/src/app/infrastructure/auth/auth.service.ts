import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable, tap } from 'rxjs';
import { HttpClient } from '@angular/common/http';
import { Router } from '@angular/router';
import { TokenStorage } from './jwt/token.service';
import { environment } from 'src/env/environment';
import { JwtHelperService } from '@auth0/angular-jwt';
import { Login } from './model/login.model';
import { AuthenticationResponse } from './model/authentication-response.model';
import { User } from './model/user.model';
import { Registration } from './model/registration.model';

@Injectable({
  providedIn: 'root'
})
export class AuthService {
  user$ = new BehaviorSubject<User>({username: "", id: 0, role: "" });

  constructor(private http: HttpClient,
    private tokenStorage: TokenStorage,
    private router: Router) { }

  login(login: Login): Observable<AuthenticationResponse> {
    return this.http
      .post<AuthenticationResponse>(`${environment.authApiHost}/login`, login)
      .pipe(
        tap((authenticationResponse) => {
          this.tokenStorage.saveAccessToken(authenticationResponse.accessToken);
          this.setUser();
        })
      );
  }

  register(registration: Registration): Observable<AuthenticationResponse> {
    return this.http
    .post<AuthenticationResponse>(`${environment.authApiHost}/register`, registration)
    .pipe(
      tap((authenticationResponse) => {
        this.tokenStorage.saveAccessToken(authenticationResponse.accessToken);
        this.setUser();
      })
    );
  }

  logout(): void {
    this.router.navigate(['/home']).then(_ => {
      this.tokenStorage.clear();
      this.user$.next({username: "", id: 0, role: "" });
      }
    );
  }

  checkIfUserExists(): void {
    const accessToken = this.tokenStorage.getAccessToken();
    if (accessToken == null) {
      return;
    }
    this.setUser();
  }

  private setUser(): void {
    const jwtHelperService = new JwtHelperService();
    const accessToken = this.tokenStorage.getAccessToken() || "";
    const decodedToken = jwtHelperService.decodeToken(accessToken);
    const user: User = {
      id: decodedToken.id,
      username: decodedToken.username,
      role: decodedToken.role,
    };
    this.user$.next(user);
  }
  getToken(): string | null {
    return this.tokenStorage.getAccessToken();
  }
}
