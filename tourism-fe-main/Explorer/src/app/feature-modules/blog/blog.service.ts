import { Injectable } from '@angular/core';
import { TokenStorage } from 'src/app/infrastructure/auth/jwt/token.service';
import { Blog, CreateBlogPayload } from './model/blog.model';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs';
import { AddCommentPayload, BlogComment, UpdateBlogPayload, UpdateCommentPayload} from './model/blog.model';

@Injectable({
  providedIn: 'root'
})
export class BlogService {

  private apiUrl = 'http://localhost:8080/api/v1/blogs';
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

  addComment(blogId: string, payload: AddCommentPayload): Observable<BlogComment> {
        const headers = this.createAuthHeaders();
        // PUTANJA: /api/v1/blogs/{id}/comments
        const url = `${this.apiUrl}/${blogId}/comments`; 
        
        return this.http.post<BlogComment>(url, payload, { headers: headers });
    }

    toggleLike(blogId: string): Observable<{ message: string }> {
        const headers = this.createAuthHeaders();
        // PUTANJA: /api/v1/blogs/{id}/like
        const url = `${this.apiUrl}/${blogId}/like`;
        
        // POST zahtev bez body-ja
        return this.http.post<{ message: string }>(url, {}, { headers: headers }); 
    }
  
getAllBlogs(): Observable<Blog[]> {
  const headers = this.createAuthHeaders();
  return this.http.get<Blog[]>(this.apiUrl, { headers });
}

getBlogById(id: string): Observable<Blog> {
  const headers = this.createAuthHeaders();
  return this.http.get<Blog>(`${this.apiUrl}/${id}`, { headers });
}

updateBlog(blogId: string, payload: UpdateBlogPayload): Observable<Blog> {
        const headers = this.createAuthHeaders();
        const url = `${this.apiUrl}/${blogId}`;

        return this.http.put<Blog>(url, payload, { headers });
    }

    updateComment(blogId: string, commentId: string, payload: UpdateCommentPayload): Observable<BlogComment> {
        const headers = this.createAuthHeaders();
        const url = `${this.apiUrl}/${blogId}/comments/${commentId}`;

        return this.http.put<BlogComment>(url, payload, { headers });
    }
}
