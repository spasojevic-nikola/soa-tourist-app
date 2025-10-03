import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { MatDialog } from '@angular/material/dialog';
import { TourService } from '../tour.service';
import { ReviewService } from '../review.service';
import { Tour } from '../model/tour.model';
import { Review, ReviewStats } from '../model/review.model';
import { MapService } from '../services/map-service.service';
import { ReviewDialogComponent } from '../review-dialog/review-dialog.component';

@Component({
  selector: 'xp-tour-details',
  templateUrl: './tour-details.component.html',
  styleUrls: ['./tour-details.component.css']
})
export class TourDetailsComponent implements OnInit {
  tour: Tour | null = null;
  isLoading = true;
  keypointAddress: string = '';
  reviews: Review[] = [];
  reviewStats: ReviewStats | null = null;
  loadingReviews = false;

  constructor(
    private route: ActivatedRoute,
    private tourService: TourService,
    private reviewService: ReviewService,
    private mapService: MapService,
    private dialog: MatDialog
  ) {}

  ngOnInit(): void {
    const tourId = this.route.snapshot.params['id'];
    this.loadTourDetails(tourId);
    this.loadReviews(tourId);
    this.loadReviewStats(tourId);
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

  loadReviews(tourId: number): void {
    this.loadingReviews = true;
    this.reviewService.getReviewsByTour(tourId).subscribe({
      next: (reviews) => {
        this.reviews = reviews;
        this.loadingReviews = false;
      },
      error: (err) => {
        console.error('Error loading reviews:', err);
        this.loadingReviews = false;
      }
    });
  }

  loadReviewStats(tourId: number): void {
    this.reviewService.getTourRatingStats(tourId).subscribe({
      next: (stats) => {
        this.reviewStats = stats;
      },
      error: (err) => {
        console.error('Error loading review stats:', err);
      }
    });
  }

  onPurchase(): void {
    // Placeholder for purchase functionality
    alert('Purchase functionality coming soon!');
  }

  onLeaveReview(): void {
    if (!this.tour) return;

    const dialogRef = this.dialog.open(ReviewDialogComponent, {
      width: '600px',
      data: {
        tourId: this.tour.id,
        tourName: this.tour.name
      }
    });

    dialogRef.afterClosed().subscribe(result => {
      if (result) {
        this.reviewService.createReview(result).subscribe({
          next: (review) => {
            console.log('Review created:', review);
            // Reload reviews and stats
            this.loadReviews(this.tour!.id);
            this.loadReviewStats(this.tour!.id);
          },
          error: (err) => {
            console.error('Error creating review:', err);
            alert('Failed to create review. Please try again.');
          }
        });
      }
    });
  }

  getRatingStars(rating: number): string[] {
    const stars = [];
    for (let i = 1; i <= 5; i++) {
      stars.push(i <= rating ? 'star' : 'star_border');
    }
    return stars;
  }
}
