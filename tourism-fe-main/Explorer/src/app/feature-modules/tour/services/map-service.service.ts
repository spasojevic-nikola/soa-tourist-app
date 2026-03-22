import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { firstValueFrom } from 'rxjs';
import { environment } from 'src/env/environment';

interface ReverseGeocodeResponse {
  city?: string;
  locality?: string;
  principalSubdivision?: string;
  countryName?: string;
}

@Injectable({
  providedIn: 'root'
})
export class MapService {
  private readonly locationCache = new Map<string, string>();
  
  constructor(private http: HttpClient) {}

  async reverseGeocode(lat: number, lng: number): Promise<string> {
    const cacheKey = this.buildCacheKey(lat, lng);
    const cached = this.locationCache.get(cacheKey);
    if (cached) {
      return cached;
    }

    const params = {
      latitude: lat.toString(),
      longitude: lng.toString(),
      localityLanguage: 'en'
    };

    try {
      const response = await firstValueFrom(
        this.http.get<ReverseGeocodeResponse>(environment.reverseGeocodeApiHost, { params })
      );
      const formatted = this.formatAddress(response);
      this.locationCache.set(cacheKey, formatted);
      return formatted;
    } catch (error) {
      console.error('Reverse geocoding failed:', error);
      return 'Address not available';
    }
  }

  private buildCacheKey(lat: number, lng: number): string {
    return `${lat.toFixed(5)},${lng.toFixed(5)}`;
  }

  private formatAddress(data: ReverseGeocodeResponse): string {
    const addressParts: string[] = [];

    const locality = data.locality || data.city;
    if (locality) {
      addressParts.push(locality);
    } else if (data.principalSubdivision) {
      addressParts.push(data.principalSubdivision);
    }

    if (data.countryName) {
      addressParts.push(data.countryName);
    }

    return addressParts.length ? addressParts.join(', ') : 'Address not found';
  }
}