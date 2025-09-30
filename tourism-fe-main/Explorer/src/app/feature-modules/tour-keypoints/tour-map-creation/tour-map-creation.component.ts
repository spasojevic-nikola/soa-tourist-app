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

        this.keyPoints.push(result);
        this.drawExistingKeyPoints();
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
    this.keyPoints.splice(index, 1);
    this.drawExistingKeyPoints();
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

}