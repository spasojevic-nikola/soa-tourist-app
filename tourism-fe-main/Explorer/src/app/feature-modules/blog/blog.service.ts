import { Injectable } from '@angular/core';
import { TokenStorage } from 'src/app/infrastructure/auth/jwt/token.service';
import { Blog, CreateBlogPayload } from './model/blog.model';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class BlogService {

  private apiUrl = 'http://localhost:8081/api/v1/blogs';
  constructor(
    private http: HttpClient, 
    private tokenStorage: TokenStorage
  ) { }
  createBlog(payload: CreateBlogPayload): Observable<Blog> {
    const headers = this.createAuthHeaders();
    
    return this.http.post<Blog>(this.apiUrl, payload, { headers: headers });
  }

  private createAuthHeaders(): HttpHeaders {
    const token = this.tokenStorage.getAccessToken();

    if (!token) {
      console.error('No token found!');
      return new HttpHeaders();
    }
    
    return new HttpHeaders({
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`
    });
  }

}
