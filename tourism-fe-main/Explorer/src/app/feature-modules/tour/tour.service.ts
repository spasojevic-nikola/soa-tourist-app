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
}
