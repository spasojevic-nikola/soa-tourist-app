import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { MatDialog } from '@angular/material/dialog';
import { TourService } from '../tour.service';
import { ReviewService } from '../review.service';
import { Tour } from '../model/tour.model';
import { Review, ReviewStats } from '../model/review.model';
import { MapService } from '../services/map-service.service';
import { ReviewDialogComponent } from '../review-dialog/review-dialog.component';
import { CartService } from '../../shopping-cart/services/cart.service';
import { MatSnackBar } from '@angular/material/snack-bar';
import { CartStateService } from '../../shopping-cart/services/cart-state.service';

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

  isAddingToCart: boolean = false; 

  constructor(
    private route: ActivatedRoute,
    private tourService: TourService,
    private reviewService: ReviewService,
    private mapService: MapService,
    private dialog: MatDialog,
    private router: Router,
    private cartService: CartService, 
    private cartStateService: CartStateService,
    private snackBar: MatSnackBar
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
    // 1. Provera postojanja podataka i uslova za kupovinu
    if (this.isAddingToCart) return; 
    const isPublished = this.tour?.status === 'Published';

    if (!this.tour || !this.tour.id || !this.tour.name) {
        this.snackBar.open('Tour details are incomplete or still loading.', 'Dismiss', { duration: 3000 });
        return;
    }
    
    if (!isPublished ) {
        let message = '';
        if (!isPublished) {
             message = 'Cannot purchase: The tour must be published.';
        } 
        this.snackBar.open(message, 'Dismiss', { duration: 4000 });
        return;
    }

    this.isAddingToCart = true; 
 
    const itemToAdd = {
        tourId: String(this.tour.id), // KONVERTOVAN NUMBER U STRING
        name: this.tour.name,
        price: this.tour.price
    };

    // 3. POZIV BACKENDA
    this.cartService.addItem(itemToAdd).subscribe({
        next: (updatedCart) => {
            this.snackBar.open(`"${itemToAdd.name}" added to cart! Total: ${updatedCart.total} RSD`, 'View Cart', { duration: 4000 })
                .onAction()
                .subscribe(() => {
                    this.router.navigate(['/shopping-cart']); 
                });
            
            // AZURIRAJ NAVBAR
            this.cartStateService.updateCartCount(updatedCart.items.length);
            this.isAddingToCart = false; 
        },
        error: (err) => {
            console.error('Error adding item to cart:', err);
            this.snackBar.open('Failed to add item. Check log for details.', 'Dismiss', { duration: 5000 });
            this.isAddingToCart = false;
        }
    });
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
