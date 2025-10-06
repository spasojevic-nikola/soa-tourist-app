import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';
import { CreateReviewRequest, Review, ReviewStats, UpdateReviewRequest } from './model/review.model';
import { environment } from 'src/env/environment';

@Injectable({
  providedIn: 'root'
})
export class ReviewService {
  private baseUrl = environment.tourApiHost;

  constructor(private http: HttpClient) {}

  private getAuthHeaders(): HttpHeaders {
    const token = localStorage.getItem('jwt');
    return new HttpHeaders({
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`
    });
  }

  // Create a new review
  createReview(request: CreateReviewRequest): Observable<Review> {
    return this.http.post<any>(
      `${this.baseUrl}/${request.tourId}/reviews`,
      request,
      { headers: this.getAuthHeaders() }
    ).pipe(
      map(review => this.parseReview(review))
    );
  }

  // Get all reviews for a tour
  getReviewsByTour(tourId: number): Observable<Review[]> {
    return this.http.get<any[]>(`${this.baseUrl}/${tourId}/reviews`).pipe(
      map(reviews => reviews.map(review => this.parseReview(review)))
    );
  }

  // Get rating statistics for a tour
  getTourRatingStats(tourId: number): Observable<ReviewStats> {
    return this.http.get<ReviewStats>(`${this.baseUrl}/${tourId}/reviews/stats`);
  }

  // Helper method to parse review images from JSON string to array
  private parseReview(review: any): Review {
    return {
      ...review,
      images: this.parseImages(review.images)
    };
  }

  private parseImages(imagesJson: string | string[]): string[] {
    // If already an array, return it
    if (Array.isArray(imagesJson)) {
      return imagesJson;
    }
    
    // If empty or null, return empty array
    if (!imagesJson || imagesJson.trim() === '') {
      return [];
    }
    
    // Try to parse JSON string
    try {
      const parsed = JSON.parse(imagesJson);
      return Array.isArray(parsed) ? parsed : [];
    } catch (e) {
      console.error('Failed to parse images JSON:', e);
      return [];
    }
  }

  // Update a review
  updateReview(reviewId: number, request: UpdateReviewRequest): Observable<Review> {
    return this.http.put<any>(
      `${this.baseUrl}/reviews/${reviewId}`,
      request,
      { headers: this.getAuthHeaders() }
    ).pipe(
      map(review => this.parseReview(review))
    );
  }

  // Delete a review
  deleteReview(reviewId: number): Observable<void> {
    return this.http.delete<void>(
      `${this.baseUrl}/reviews/${reviewId}`,
      { headers: this.getAuthHeaders() }
    );
  }

  // Get all reviews by the authenticated user
  getMyReviews(): Observable<Review[]> {
    return this.http.get<any[]>(
      `${this.baseUrl}/my-reviews`,
      { headers: this.getAuthHeaders() }
    ).pipe(
      map(reviews => reviews.map(review => this.parseReview(review)))
    );
  }
}
