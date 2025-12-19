import { Component, OnInit } from '@angular/core';
import { Observable } from 'rxjs';
import { TourService } from '../tour.service';
import { Tour } from '../model/tour.model';

@Component({
  selector: 'xp-tour-list',
  templateUrl: './tour-list.component.html',
  styleUrls: ['./tour-list.component.css']
})
export class TourListComponent implements OnInit {
  tours$: Observable<Tour[]>;

  constructor(private tourService: TourService) { }

  ngOnInit(): void {
    this.loadTours();
  }

  loadTours(): void {
    this.tours$ = this.tourService.getAuthorTours();
  }

  publishTour(tourId: number): void {
    this.tourService.publishTour(tourId).subscribe({
      next: () => {
        alert('Tour successfully published!');
        this.loadTours();
      },
      error: (err) => {
        alert(`Failed to publish tour: ${err.error || err.message}`);
      }
    });
  }

  archiveTour(tourId: number): void {
    if (confirm('Are you sure you want to archive this tour?')) {
      this.tourService.archiveTour(tourId).subscribe({
        next: () => {
          alert('Tour successfully archived!');
          this.loadTours();
        },
        error: (err) => {
          alert(`Failed to archive tour: ${err.error || err.message}`);
        }
      });
    }
  }

  activateTour(tourId: number): void {
    this.tourService.activateTour(tourId).subscribe({
      next: () => {
        alert('Tour successfully reactivated!');
        this.loadTours();
      },
      error: (err) => {
        alert(`Failed to activate tour: ${err.error || err.message}`);
      }
    });
  }
}