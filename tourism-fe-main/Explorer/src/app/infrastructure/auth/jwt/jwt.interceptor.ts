import { HttpEvent, HttpHandler, HttpInterceptor, HttpRequest } from "@angular/common/http";
import { Injectable } from "@angular/core";
import { Observable } from "rxjs";
import { ACCESS_TOKEN } from '../../../shared/constants';
import { JwtHelperService } from "@auth0/angular-jwt";

@Injectable()
export class JwtInterceptor implements HttpInterceptor {
  private jwtHelper = new JwtHelperService();

  intercept(request: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {
    const token = localStorage.getItem(ACCESS_TOKEN);
    let headers: any = {};

    if (token) {
      headers['Authorization'] = `Bearer ${token}`;

      // Dekoduj token da bi izvukao ID korisnika
      const decodedToken = this.jwtHelper.decodeToken(token);
      if (decodedToken && decodedToken.id) {
        headers['X-User-ID'] = decodedToken.id.toString();
      }
    }

    const authReq = request.clone({ setHeaders: headers });
    return next.handle(authReq);
  }
}