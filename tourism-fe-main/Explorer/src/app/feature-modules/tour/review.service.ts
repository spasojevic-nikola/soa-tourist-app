import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { JwtHelperService } from '@auth0/angular-jwt';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';
import { TokenStorage } from 'src/app/infrastructure/auth/jwt/token.service';
import { environment } from 'src/env/environment';
import {
  CreateReviewRequest,
  Review,
  ReviewStats,
  UpdateReviewRequest
} from './model/review.model';

@Injectable({
  providedIn: 'root'
})
export class ReviewService {
  private readonly baseUrl = environment.tourApiHost;
  private readonly jwtHelper = new JwtHelperService();

  constructor(private http: HttpClient, private tokenStorage: TokenStorage) {}

  // Ensure the backend always receives the authenticated user headers expected by tour-service.
  private getAuthHeaders(): HttpHeaders {
    const token = this.tokenStorage.getAccessToken();
    const headers: Record<string, string> = {
      'Content-Type': 'application/json'
    };

    if (token) {
      headers['Authorization'] = `Bearer ${token}`;

      try {
        const decoded = this.jwtHelper.decodeToken(token) as {
          id?: number;
          username?: string;
          role?: string;
        };

        if (decoded?.id) {
          headers['X-User-ID'] = decoded.id.toString();
        }
        if (decoded?.username) {
          headers['X-User-Username'] = decoded.username;
        }
        if (decoded?.role) {
          headers['X-User-Role'] = decoded.role;
        }
      } catch (error) {
        console.warn('Unable to decode JWT for review headers', error);
      }
    }

    return new HttpHeaders(headers);
  }

  createReview(request: CreateReviewRequest): Observable<Review> {
    return this.http
      .post<any>(`${this.baseUrl}/${request.tourId}/reviews`, request, {
        headers: this.getAuthHeaders()
      })
      .pipe(map((review) => this.parseReview(review)));
  }

  getReviewsByTour(tourId: number): Observable<Review[]> {
    return this.http
      .get<any[]>(`${this.baseUrl}/${tourId}/reviews`)
      .pipe(map((reviews) => reviews.map((review) => this.parseReview(review))));
  }

  getTourRatingStats(tourId: number): Observable<ReviewStats> {
    return this.http.get<ReviewStats>(`${this.baseUrl}/${tourId}/reviews/stats`);
  }

  updateReview(reviewId: number, request: UpdateReviewRequest): Observable<Review> {
    return this.http
      .put<any>(`${this.baseUrl}/reviews/${reviewId}`, request, {
        headers: this.getAuthHeaders()
      })
      .pipe(map((review) => this.parseReview(review)));
  }

  deleteReview(reviewId: number): Observable<void> {
    return this.http.delete<void>(`${this.baseUrl}/reviews/${reviewId}`, {
      headers: this.getAuthHeaders()
    });
  }

  getMyReviews(): Observable<Review[]> {
    return this.http
      .get<any[]>(`${this.baseUrl}/my-reviews`, {
        headers: this.getAuthHeaders()
      })
      .pipe(map((reviews) => reviews.map((review) => this.parseReview(review))));
  }

  private parseReview(review: any): Review {
    return {
      ...review,
      images: this.parseImages(review.images)
    };
  }

  private parseImages(images: string | string[]): string[] {
    if (Array.isArray(images)) {
      return images;
    }

    if (!images || !images.trim()) {
      return [];
    }

    try {
      const parsed = JSON.parse(images);
      return Array.isArray(parsed) ? parsed : [];
    } catch (error) {
      console.error('Failed to parse review images JSON', error);
      return [];
    }
  }
}
