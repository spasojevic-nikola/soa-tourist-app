import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { TourService } from '../tour.service';
import { Tour } from '../model/tour.model';
import { MapService } from '../services/map-service.service';

@Component({
  selector: 'xp-tour-details',
  templateUrl: './tour-details.component.html',
  styleUrls: ['./tour-details.component.css']
})
export class TourDetailsComponent implements OnInit {
  tour: Tour | null = null;
  isLoading = true;
  keypointAddress: string = '';

  constructor(
    private route: ActivatedRoute,
    private tourService: TourService,
    private mapService: MapService
  ) {}

  ngOnInit(): void {
    const tourId = this.route.snapshot.params['id'];
    this.loadTourDetails(tourId);
  }

  loadTourDetails(tourId: number): void {
    this.tourService.getTourById(tourId).subscribe({
      next: async (tour) => {
        this.tour = tour;
        // Load address for first keypoint if exists
        if (tour.keyPoints && tour.keyPoints.length > 0) {
          const firstKeypoint = tour.keyPoints[0];
          if (firstKeypoint.address) {
            this.keypointAddress = firstKeypoint.address;
          } else {
            this.keypointAddress = await this.mapService.reverseGeocode(
              firstKeypoint.latitude,
              firstKeypoint.longitude
            );
          }
        }
        this.isLoading = false;
      },
      error: (err) => {
        console.error('Error loading tour details:', err);
        this.isLoading = false;
      }
    });
  }

  onPurchase(): void {
    // Placeholder for purchase functionality
    alert('Purchase functionality coming soon!');
  }

  onLeaveReview(): void {
    // Placeholder for review functionality
    alert('Review functionality coming soon!');
  }
}
