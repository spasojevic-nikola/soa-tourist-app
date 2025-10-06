import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { environment } from 'src/env/environment';
import { Tour } from './model/tour.model';
import { Observable } from 'rxjs';
import { CreateTourPayload } from './dto/tour-creation.dto';

@Injectable({
  providedIn: 'root'
})
export class TourService {
  private apiUrl = environment.tourApiHost;


  constructor(private http: HttpClient) { }
  
  createTour(payload: CreateTourPayload): Observable<Tour> {
    return this.http.post<Tour>(`${this.apiUrl}/create-tour`, payload);
  }
  
  getAuthorTours(): Observable<Tour[]> {
    return this.http.get<Tour[]>(this.apiUrl);
  }

  getAllPublishedTours(): Observable<Tour[]> {
    return this.http.get<Tour[]>(`${this.apiUrl}/published`);
  }

  getTourById(tourId: number): Observable<Tour> {
    return this.http.get<Tour>(`${this.apiUrl}/${tourId}`);
  }

  addDuration(tourId: number, payload: any): Observable<any> {
    return this.http.post(`${this.apiUrl}/${tourId}/duration`, payload);
  }

  publishTour(tourId: number): Observable<Tour> {
    return this.http.put<Tour>(`${this.apiUrl}/${tourId}/publish`, {});
  }

  archiveTour(tourId: number): Observable<Tour> {
    return this.http.put<Tour>(`${this.apiUrl}/${tourId}/archive`, {});
  }

  activateTour(tourId: number): Observable<Tour> {
    return this.http.put<Tour>(`${this.apiUrl}/${tourId}/activate`, {});
  }
}
