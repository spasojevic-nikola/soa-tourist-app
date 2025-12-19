import { Component, OnInit, AfterViewInit, Output, EventEmitter } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { CreateKeyPointPayload } from '../model/keypoint.model';
import { KeypointDialogService } from '../../tour/services/keypoint-dialog.service';
import { MapService } from '../../tour/services/map-service.service';
import { KeyPoint } from '../model/keypoint.model';
import { Observable, of, from } from 'rxjs';
import { tap, catchError } from 'rxjs/operators';

declare let L: any;

@Component({
  selector: 'xp-tour-map-creation',
  templateUrl: './tour-map-creation.component.html',
  styleUrls: ['./tour-map-creation.component.css']
})
export class TourMapCreationComponent implements OnInit, AfterViewInit {
  @Output() keyPointsCompleted = new EventEmitter<CreateKeyPointPayload[]>();
  @Output() goBack = new EventEmitter<void>();

  map: any;
  keyPoints: CreateKeyPointPayload[] = [];
  selectedPosition: { lat: number; lng: number } | null = null;
  markers: any[] = [];
  polyline: any;
  tempMarker: any;
  private addressCache = new Map<string, string>();
  private imagePreviews = new Map<number, string>();
  isEditMode = false;
  editingIndex: number | null = null;
  totalDistance: number = 0;

  constructor(
    private fb: FormBuilder,
    private keypointDialogService: KeypointDialogService,
    private mapService: MapService
  ) {}

  ngOnInit(): void {}

  ngAfterViewInit(): void {
    setTimeout(() => {
      this.initMap();
    }, 100);
  }

  initMap(): void {
    this.map = L.map('map').setView([44.7866, 20.4489], 13);

    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
      attribution: '© OpenStreetMap contributors'
    }).addTo(this.map);

    this.map.on('click', (e: any) => {
      this.onMapClick(e);
    });
  }

  onMapClick(e: any): void {
    if (this.isEditMode && this.editingIndex !== null) {
      this.finishEditKeyPoint(e.latlng.lat, e.latlng.lng);
      return;
    }
    
    if (this.tempMarker) {
      this.map.removeLayer(this.tempMarker);
    }

    this.tempMarker = L.marker([e.latlng.lat, e.latlng.lng], {
      icon: L.divIcon({
        className: 'temp-marker',
        html: '📍',
        iconSize: [30, 30]
      })
    }).addTo(this.map);

    this.selectedPosition = { lat: e.latlng.lat, lng: e.latlng.lng };
    this.openKeypointDialog(e.latlng.lat, e.latlng.lng);
  }

  openKeypointDialog(lat: number, lng: number): void {
    const order = this.keyPoints.length + 1;
    
    this.keypointDialogService.openKeypointDialog(lat, lng, order).subscribe(async (result) => {
      if (result) {
        try {
          result.address = await this.mapService.reverseGeocode(result.latitude, result.longitude);
        } catch (error) {
          result.address = `Lat: ${result.latitude.toFixed(4)}, Lng: ${result.longitude.toFixed(4)}`;
        }

        if (result.image && typeof result.image === 'string') {
          this.imagePreviews.set(this.keyPoints.length, result.image);
        }

        this.keyPoints.push(result);
        this.drawExistingKeyPoints();
        this.calculateTotalDistance();
      }
      
      if (this.tempMarker) {
        this.map.removeLayer(this.tempMarker);
        this.tempMarker = null;
      }
    });
  }

  drawExistingKeyPoints(): void {
    this.markers.forEach(marker => this.map.removeLayer(marker));
    if (this.polyline) {
      this.map.removeLayer(this.polyline);
    }

    this.markers = [];
    const latlngs: [number, number][] = [];

    const sortedKeyPoints = [...this.keyPoints].sort((a, b) => a.order - b.order);

    sortedKeyPoints.forEach((kp, index) => {
      const marker = L.marker([kp.latitude, kp.longitude])
        .addTo(this.map)
        .bindPopup(`
          <div class="popup-content">
            <strong>${kp.order}. ${kp.name}</strong><br>
            ${kp.description}
          </div>
        `);

      this.markers.push(marker);
      latlngs.push([kp.latitude, kp.longitude]);
    });

    if (latlngs.length > 1) {
      this.polyline = L.polyline(latlngs, { color: 'blue', weight: 4 }).addTo(this.map);
    }
  }

  removeKeyPoint(index: number): void {
    this.imagePreviews.delete(index);
    this.keyPoints.splice(index, 1);
    
    const newImagePreviews = new Map<number, string>();
    this.keyPoints.forEach((kp, newIndex) => {
      const oldPreview = this.imagePreviews.get(newIndex + 1);
      if (oldPreview) {
        newImagePreviews.set(newIndex, oldPreview);
      }
    });
    this.imagePreviews = newImagePreviews;
    
    this.drawExistingKeyPoints();
    this.calculateTotalDistance();
  }

  finishTour(): void {
    if (this.keyPoints.length < 2) {
      alert('Please add at least two key points before creating the tour.');
      return;
    }
    this.keyPointsCompleted.emit(this.keyPoints);
  }

  backToStep1(): void {
    this.goBack.emit();
  }

  getAddressDisplay(keyPoint: CreateKeyPointPayload): Observable<string> {
    const cacheKey = `${keyPoint.latitude},${keyPoint.longitude}`;
    
    if (this.addressCache.has(cacheKey)) {
      return of(this.addressCache.get(cacheKey)!);
    }
    
    return from(this.mapService.reverseGeocode(keyPoint.latitude, keyPoint.longitude)).pipe(
      tap(address => this.addressCache.set(cacheKey, address)),
      catchError(error => {
        const fallback = `Lat: ${keyPoint.latitude.toFixed(4)}, Lng: ${keyPoint.longitude.toFixed(4)}`;
        this.addressCache.set(cacheKey, fallback);
        return of(fallback);
      })
    );
  }

  getImagePreview(index: number): string | null {
    return this.imagePreviews.get(index) || null;
  }

  startEditKeyPoint(index: number): void {
    this.isEditMode = true;
    this.editingIndex = index;
  }

  async finishEditKeyPoint(lat: number, lng: number): Promise<void> {
    if (this.editingIndex === null) return;

    const existingKeyPoint = this.keyPoints[this.editingIndex];
    
    // Dodaj temporary marker na novu lokaciju
    if (this.tempMarker) {
      this.map.removeLayer(this.tempMarker);
    }
    
    this.tempMarker = L.marker([lat, lng], {
      icon: L.divIcon({
        className: 'temp-marker',
        html: '📍',
        iconSize: [30, 30]
      })
    }).addTo(this.map);

    // Kreiraj ažurirani key point sa NOVOM LOKACIJOM
    const updatedKeyPoint: CreateKeyPointPayload = {
      ...existingKeyPoint,
      latitude: lat,  // Koristi NOVE koordinate
      longitude: lng  // Koristi NOVE koordinate
    };

    // Otvori edit dijalog sa ažuriranim podacima
    this.openEditKeypointDialog(lat, lng, this.editingIndex, updatedKeyPoint);
  }

  cancelEdit(): void {
    this.isEditMode = false;
    this.editingIndex = null;
    
    if (this.tempMarker) {
      this.map.removeLayer(this.tempMarker);
      this.tempMarker = null;
    }
  }

  openEditKeypointDialog(lat: number, lng: number, index: number, existingKeyPoint: CreateKeyPointPayload): void {
    this.keypointDialogService.openKeypointDialog(
      lat, 
      lng, 
      existingKeyPoint.order,
      existingKeyPoint
    ).subscribe(async (result) => {
      if (result) {
        try {
          result.address = await this.mapService.reverseGeocode(result.latitude, result.longitude);
        } catch (error) {
          result.address = `Lat: ${result.latitude.toFixed(4)}, Lng: ${result.longitude.toFixed(4)}`;
        }

        // Ažuriraj postojeći key point
        this.keyPoints[index] = result;
        
        if (result.image && typeof result.image === 'string') {
          this.imagePreviews.set(index, result.image);
        }
        
        this.drawExistingKeyPoints();
        this.calculateTotalDistance();
      }
      
      this.isEditMode = false;
      this.editingIndex = null;
      
      if (this.tempMarker) {
        this.map.removeLayer(this.tempMarker);
        this.tempMarker = null;
      }
    });
  }

  /**
   * Calculates total distance between all keypoints using Haversine formula
   */
  calculateTotalDistance(): void {
    if (this.keyPoints.length < 2) {
      this.totalDistance = 0;
      return;
    }

    let distance = 0;
    const sortedKeyPoints = [...this.keyPoints].sort((a, b) => a.order - b.order);

    for (let i = 0; i < sortedKeyPoints.length - 1; i++) {
      const kp1 = sortedKeyPoints[i];
      const kp2 = sortedKeyPoints[i + 1];
      distance += this.haversineDistance(
        kp1.latitude,
        kp1.longitude,
        kp2.latitude,
        kp2.longitude
      );
    }

    this.totalDistance = distance;
  }

  /**
   * Haversine formula - calculates distance between two GPS coordinates in kilometers
   */
  private haversineDistance(lat1: number, lon1: number, lat2: number, lon2: number): number {
    const earthRadius = 6371; // Earth radius in kilometers

    // Convert degrees to radians
    const lat1Rad = (lat1 * Math.PI) / 180;
    const lat2Rad = (lat2 * Math.PI) / 180;
    const deltaLat = ((lat2 - lat1) * Math.PI) / 180;
    const deltaLon = ((lon2 - lon1) * Math.PI) / 180;

    // Haversine formula
    const a =
      Math.sin(deltaLat / 2) * Math.sin(deltaLat / 2) +
      Math.cos(lat1Rad) * Math.cos(lat2Rad) * Math.sin(deltaLon / 2) * Math.sin(deltaLon / 2);
    const c = 2 * Math.asin(Math.sqrt(a));

    return earthRadius * c;
  }

  /**
   * Calculate estimated time by car (average speed: 40 km/h in city)
   */
  getCarTime(): number {
    if (this.totalDistance === 0) return 0;
    return (this.totalDistance / 40) * 60; // returns minutes
  }

  /**
   * Calculate estimated time by bicycle (average speed: 15 km/h)
   */
  getBicycleTime(): number {
    if (this.totalDistance === 0) return 0;
    return (this.totalDistance / 15) * 60; // returns minutes
  }

  /**
   * Calculate estimated time walking (average speed: 5 km/h)
   */
  getWalkingTime(): number {
    if (this.totalDistance === 0) return 0;
    return (this.totalDistance / 5) * 60; // returns minutes
  }
}