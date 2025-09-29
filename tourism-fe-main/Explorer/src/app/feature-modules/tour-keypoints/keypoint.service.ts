import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from 'src/env/environment';
import { KeyPoint, CreateKeyPointPayload } from './model/keypoint.model';

@Injectable({
  providedIn: 'root'
})
export class KeypointService {
  private apiUrl = environment.tourApiHost;

  constructor(private http: HttpClient) { }

  createKeyPoint(tourId: number, payload: CreateKeyPointPayload): Observable<KeyPoint> {
    return this.http.post<KeyPoint>(`${this.apiUrl}/${tourId}/keypoints`, payload);
  }

  getKeyPointsByTour(tourId: number): Observable<KeyPoint[]> {
    return this.http.get<KeyPoint[]>(`${this.apiUrl}/${tourId}/keypoints`);
  }

  updateKeyPoint(keyPointId: number, payload: CreateKeyPointPayload): Observable<KeyPoint> {
    return this.http.put<KeyPoint>(`${this.apiUrl}/keypoints/${keyPointId}`, payload);
  }

  deleteKeyPoint(keyPointId: number): Observable<void> {
    return this.http.delete<void>(`${this.apiUrl}/keypoints/${keyPointId}`);
  }
}