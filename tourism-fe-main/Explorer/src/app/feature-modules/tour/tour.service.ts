import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { environment } from 'src/env/environment';
import { Tour } from './model/tour.model';
import { Observable } from 'rxjs';
import { CreateTourPayload, UpdateTourPayload } from './dto/tour-creation.dto';

@Injectable({
  providedIn: 'root'
})
export class TourService {
  private apiUrl = environment.tourApiHost;


  constructor(private http: HttpClient) { }
  
  createTour(payload: CreateTourPayload): Observable<Tour> {
    return this.http.post<Tour>(`${this.apiUrl}/create-tour`, payload);
  }

updateTour(id: number, payload: UpdateTourPayload): Observable<Tour> {
    return this.http.put<Tour>(`${this.apiUrl}/${id}`, payload);
  }

  getAuthorTours(): Observable<Tour[]> {
    return this.http.get<Tour[]>(this.apiUrl);
  }

  getPublishedTours(): Observable<Tour[]> {
    return this.http.get<Tour[]>(`${this.apiUrl}/published`);
  }

  getTourById(tourId: number): Observable<Tour> {
    return this.http.get<Tour>(`${this.apiUrl}/${tourId}`);
  }

  getAllExecutionsForTour(tourId: number): Observable<any[]> {
    return this.http.get<any[]>(`${this.apiUrl}/executions/tour/${tourId}`);
  }

  startTourExecution(tourId: number, startLat: number = 0, startLng: number = 0): Observable<any> {
    return this.http.post<any>(`${this.apiUrl}/${tourId}/start`, {
      startLat,
      startLng
    });
  }
}
