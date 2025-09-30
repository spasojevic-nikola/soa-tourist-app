import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';   
import { firstValueFrom } from 'rxjs';   

@Injectable({
  providedIn: 'root'
})
export class MapService {
  
  constructor(private http: HttpClient) {}  // ← Inject HttpClient

  async reverseGeocode(lat: number, lng: number): Promise<string> {
    const url = `https://nominatim.openstreetmap.org/reverse?format=json&lat=${lat}&lon=${lng}`;
    
    try {
      // Koristi firstValueFrom umjesto toPromise()
      const response: any = await firstValueFrom(this.http.get(url));
      
      if (response && response.address) {
        const addr = response.address;
        const addressParts = [];
        
        if (addr.road) addressParts.push(addr.road);
        if (addr.house_number) addressParts.push(addr.house_number);
        if (addr.city || addr.town || addr.village) addressParts.push(addr.city || addr.town || addr.village);
        if (addr.country) addressParts.push(addr.country);
        
        return addressParts.join(', ');
      }
      
      return 'Address not found';
    } catch (error) {
      console.error('Reverse geocoding failed:', error);
      return 'Address not available';
    }
  }
}