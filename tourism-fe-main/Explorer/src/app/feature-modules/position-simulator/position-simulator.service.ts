import { Injectable } from '@angular/core';

export interface Position {
  lat: number;
  lng: number;
}

@Injectable({
  providedIn: 'root'
})
export class PositionSimulatorService {
  private readonly STORAGE_KEY = 'tourist-position';

  constructor() { }

  getCurrentPosition(): Position | null {
    const stored = localStorage.getItem(this.STORAGE_KEY);
    if (stored) {
      try {
        return JSON.parse(stored);
      } catch (error) {
        console.error('Error parsing stored position:', error);
        return null;
      }
    }
    return null;
  }

  savePosition(position: Position): void {
    localStorage.setItem(this.STORAGE_KEY, JSON.stringify(position));
  }

  clearPosition(): void {
    localStorage.removeItem(this.STORAGE_KEY);
  }

  hasPosition(): boolean {
    return this.getCurrentPosition() !== null;
  }
}