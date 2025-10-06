import { Component, OnInit, AfterViewInit } from '@angular/core';
import { PositionSimulatorService } from './position-simulator.service';
import * as L from 'leaflet';

@Component({
  selector: 'app-position-simulator',
  templateUrl: './position-simulator.component.html',
  styleUrls: ['./position-simulator.component.css']
})
export class PositionSimulatorComponent implements OnInit, AfterViewInit {
  currentPosition: { lat: number; lng: number; } | null = null;
  map: L.Map | null = null;
  currentMarker: L.Marker | null = null;
  isMapLoaded = false;

  manualLatitude: number | null = null;
  manualLongitude: number | null = null;

  constructor(private positionService: PositionSimulatorService) { }

  ngOnInit(): void {
    this.loadCurrentPosition();
  }

  ngAfterViewInit(): void {
    // Dodaj mali delay da se osigura da je DOM potpuno renderovan
    setTimeout(() => {
      this.initializeMap();
    }, 100);
  }

  initializeMap(): void {
    try {
      // Default position (Novi Sad centar)
      const defaultPosition: [number, number] = [45.2671, 19.8335];
      
      // Kreiraj mapu
      this.map = L.map('map', {
        center: defaultPosition,
        zoom: 13,
        zoomControl: true,
        scrollWheelZoom: true
      });

      // Dodaj tile layer (OpenStreetMap)
      L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
        maxZoom: 19,
        attribution: '© OpenStreetMap contributors'
      }).addTo(this.map);

      this.isMapLoaded = true;

      // Dodaj click listener na mapu
      this.map.on('click', (event: L.LeafletMouseEvent) => {
        const lat = event.latlng.lat;
        const lng = event.latlng.lng;
        this.updatePosition(lat, lng);
      });

      // Ako već postoji pozicija, prikaži je
      if (this.currentPosition) {
        this.updateMapCenter(this.currentPosition);
        this.addMarker(this.currentPosition);
      }

    } catch (error) {
      console.error('Error initializing map:', error);
      this.isMapLoaded = false;
    }
  }

  loadCurrentPosition(): void {
    const position = this.positionService.getCurrentPosition();
    if (position) {
      this.currentPosition = position;
    }
  }

  updatePosition(lat: number, lng: number): void {
    const newPosition = { lat, lng };
    
    // Sačuvaj poziciju u servisu
    this.positionService.savePosition(newPosition);
    this.currentPosition = newPosition;
    
    // Ažuriraj marker na mapi
    this.updateMarker(newPosition);
  }

  private updateMapCenter(position: { lat: number; lng: number; }): void {
    if (this.map && this.isMapLoaded) {
      this.map.setView([position.lat, position.lng], 13);
    }
  }

  private addMarker(position: { lat: number; lng: number; }): void {
    if (this.map && this.isMapLoaded) {
      // Ukloni postojeći marker
      if (this.currentMarker) {
        this.map.removeLayer(this.currentMarker);
      }

      // Kreiraj custom icon
      const customIcon = L.divIcon({
        className: 'custom-div-icon',
        html: `
          <div style="
            background-color: #4285F4;
            width: 24px;
            height: 24px;
            border-radius: 50%;
            border: 3px solid white;
            box-shadow: 0 2px 8px rgba(0,0,0,0.3);
            display: flex;
            align-items: center;
            justify-content: center;
          ">
            <div style="
              background-color: white;
              width: 8px;
              height: 8px;
              border-radius: 50%;
            "></div>
          </div>
        `,
        iconSize: [30, 30],
        iconAnchor: [15, 15]
      });

      // Dodaj marker sa custom icon-om
      this.currentMarker = L.marker([position.lat, position.lng], {
        icon: customIcon,
        draggable: true
      }).addTo(this.map);

      // Dodaj tooltip
      this.currentMarker.bindTooltip('Vaša trenutna pozicija', {
        permanent: false,
        direction: 'top'
      });

      // Dodaj dragend listener
      this.currentMarker.on('dragend', (event: L.DragEndEvent) => {
        const marker = event.target as L.Marker;
        const position = marker.getLatLng();
        this.updatePosition(position.lat, position.lng);
      });
    }
  }

  private updateMarker(position: { lat: number; lng: number; }): void {
    if (this.currentMarker && this.isMapLoaded) {
      this.currentMarker.setLatLng([position.lat, position.lng]);
    } else {
      this.addMarker(position);
    }
  }

  clearPosition(): void {
    this.positionService.clearPosition();
    this.currentPosition = null;
    if (this.currentMarker && this.map) {
      this.map.removeLayer(this.currentMarker);
      this.currentMarker = null;
    }
  }

  getCurrentLocation(): void {
    if (navigator.geolocation) {
      navigator.geolocation.getCurrentPosition(
        (position) => {
          const lat = position.coords.latitude;
          const lng = position.coords.longitude;
          this.updatePosition(lat, lng);
          this.updateMapCenter({ lat, lng });
        },
        (error) => {
          console.error('Error getting location:', error);
          alert('Ne mogu da dohvatim vašu trenutnu lokaciju. Koristite klik na mapu.');
        }
      );
    } else {
      alert('Geolocation nije podržana u vašem brauzeru.');
    }
  }

  setManualPosition(latValue: string, lngValue: string): void {
    const lat = parseFloat(latValue);
    const lng = parseFloat(lngValue);
    
    if (isNaN(lat) || isNaN(lng)) {
      alert('Molimo unesite validne numeričke vrednosti za koordinate.');
      return;
    }
    
    if (lat < -90 || lat > 90 || lng < -180 || lng > 180) {
      alert('Koordinate nisu u validnom opsegu. Latitude: -90 do 90, Longitude: -180 do 180.');
      return;
    }
    
    this.updatePosition(lat, lng);
    
    // Centriraj mapu na novu poziciju
    if (this.isMapLoaded) {
      this.updateMapCenter({ lat, lng });
    }
  }

  goToManualCoordinates(): void {
  if (this.manualLatitude !== null && this.manualLongitude !== null) {
    // Validacija koordinata
    if (this.isValidCoordinate(this.manualLatitude, this.manualLongitude)) {
      this.updatePosition(this.manualLatitude, this.manualLongitude);
      this.updateMapCenter({ lat: this.manualLatitude, lng: this.manualLongitude });
      
      this.showSuccess(`Lokacija postavljena na: ${this.manualLatitude}, ${this.manualLongitude}`);
    } else {
      this.showError('Nevalidne koordinate! Proverite unos.');
    }
  }
}

private isValidCoordinate(lat: number, lng: number): boolean {
  return lat >= -90 && lat <= 90 && lng >= -180 && lng <= 180;
}

private showSuccess(message: string): void {
  console.log('SUCCESS:', message);
  alert('✅ ' + message);
}

private showError(message: string): void {
  console.error('ERROR:', message);
  alert('❌ ' + message);
}
}