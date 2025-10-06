import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from 'src/env/environment';

@Injectable({
  providedIn: 'root'
})
export class TourExecutionService {
  private apiUrl = environment.tourApiHost;

  constructor(private http: HttpClient) { }

  getExecutionDetails(executionId: number): Observable<any> {
    return this.http.get(`${this.apiUrl}/executions/${executionId}`);
  }

  checkPosition(executionId: number, lat: number, lng: number): Observable<any> {
  const payload = {
    currentLat: lat,
    currentLng: lng
  };

  console.log('📦 Sending payload to backend:', payload);

  return this.http.post(`${this.apiUrl}/executions/${executionId}/check-position`, payload);
}

  completeTour(executionId: number): Observable<any> {
    return this.http.put(`${this.apiUrl}/executions/${executionId}/complete`, {});
  }

  abandonTour(executionId: number): Observable<any> {
    return this.http.put(`${this.apiUrl}/executions/${executionId}/abandon`, {});
  }

  getActiveExecution(tourId: number): Observable<any> {
    return this.http.get(`${this.apiUrl}/executions/active/${tourId}`);
  }

  startTourExecution(tourId: number, startLat: number, startLng: number): Observable<any> {
    return this.http.post(`${this.apiUrl}/${tourId}/start`, {
      startLat: startLat,
      startLng: startLng
    });
  }
}