import { Component, OnInit, AfterViewInit, Output, EventEmitter } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { CreateKeyPointPayload } from '../model/keypoint.model';

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
  keyPoints: CreateKeyPointPayload[] = []; // Promenjeno u CreateKeyPointPayload[]
  keyPointsForm: FormGroup;
  isFormVisible = false;
  selectedPosition: { lat: number; lng: number } | null = null;
  markers: any[] = [];
  polyline: any;
  tempMarker: any;

  constructor(
    private fb: FormBuilder
  ) {}

  ngOnInit(): void {
    this.initForm();
  }

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
    this.isFormVisible = true;
    
    this.keyPointsForm.patchValue({
      latitude: e.latlng.lat,
      longitude: e.latlng.lng,
      order: this.keyPoints.length + 1
    });
  }

  initForm(): void {
    this.keyPointsForm = this.fb.group({
      name: ['', [Validators.required, Validators.minLength(3)]],
      description: ['', [Validators.required, Validators.minLength(10)]],
      latitude: ['', Validators.required],
      longitude: ['', Validators.required],
      image: [''],
      order: [1, [Validators.required, Validators.min(1)]]
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

  onSubmit(): void {
    if (this.keyPointsForm.invalid) return;

    const payload: CreateKeyPointPayload = {
      ...this.keyPointsForm.value,
      latitude: parseFloat(this.keyPointsForm.value.latitude),
      longitude: parseFloat(this.keyPointsForm.value.longitude),
      order: parseInt(this.keyPointsForm.value.order)
    };

    // Dodaj keypoint lokalno (ne šalji na backend)
    this.keyPoints.push(payload);
    this.drawExistingKeyPoints();
    this.closeForm();
    
    if (this.tempMarker) {
      this.map.removeLayer(this.tempMarker);
      this.tempMarker = null;
    }
  }

  closeForm(): void {
    this.isFormVisible = false;
    this.selectedPosition = null;
    this.keyPointsForm.reset({ order: this.keyPoints.length + 1 });
    
    if (this.tempMarker) {
      this.map.removeLayer(this.tempMarker);
      this.tempMarker = null;
    }
  }

  removeKeyPoint(index: number): void {
    this.keyPoints.splice(index, 1);
    this.drawExistingKeyPoints();
  }

  // Nova metoda za završetak - emituje keypoints wizard-u
  finishTour(): void {
    if (this.keyPoints.length === 0) {
      alert('Please add at least one key point before creating the tour.');
      return;
    }
    this.keyPointsCompleted.emit(this.keyPoints);
  }

  // Nova metoda za povratak na step 1
  backToStep1(): void {
    this.goBack.emit();
  }
}