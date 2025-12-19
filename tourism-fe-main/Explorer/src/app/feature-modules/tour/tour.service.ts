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

  getActiveExecution(tourId: number): Observable<any> {
  return this.http.get(`${this.apiUrl}/executions/active/${tourId}`);
}

getAllExecutionsForTour(tourId: number): Observable<any[]> {
  return this.http.get<any[]>(`${this.apiUrl}/executions/tour/${tourId}`);
}

startTourExecution(tourId: number): Observable<any> {
  const position = this.getCurrentPosition(); 
  
  return this.http.post(`${this.apiUrl}/${tourId}/start`, {
    startLat: position.lat,
    startLng: position.lng
  });
}

private getCurrentPosition(): { lat: number; lng: number } {
    const storedPosition = localStorage.getItem('tourist-position');
    
    if (storedPosition) {
      try {
        return JSON.parse(storedPosition);
      } catch (e) {
        console.error('Error parsing stored position:', e);
      }
    }

    // Fallback: koristi default poziciju (Novi Sad centar)
    return {
      lat: 45.2671,
      lng: 19.8335
    };
  }
}
