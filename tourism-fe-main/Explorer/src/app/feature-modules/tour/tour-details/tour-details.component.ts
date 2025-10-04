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
import { forkJoin } from 'rxjs';

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
  hasPurchased = false;

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
    this.checkIfTourInCart(tourId); // Dodato: provera da li je tura već u korpi
  }

  checkIfTourInCart(tourId: number): void {
    this.cartService.getCart().subscribe({
      next: (cart) => {
        // Provera da li postoji item u korpi sa ovim tourId
        const tourInCart = cart.items.some(item => 
          item.tourId === String(tourId)
        );
        this.hasPurchased = tourInCart;
      },
      error: (err) => {
        console.error('Error checking cart:', err);
      }
    });
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
    if (this.isAddingToCart) return;
    
    // Prvo se proverava 'hasPurchased' status
    if (this.hasPurchased) {
      this.snackBar.open('You have already purchased this tour.', 'Dismiss', { duration: 3000 });
      return;
    }

    if (!this.tour || !this.tour.id || this.tour.status !== 'Published') {
      this.snackBar.open('This tour is not available for purchase.', 'Dismiss', { duration: 3000 });
      return;
    }

    this.isAddingToCart = true;
    const itemToAdd = { tourId: String(this.tour.id) }; // Sada se šalje samo tourId

    this.cartService.addItem(itemToAdd).subscribe({
      next: (updatedCart) => {
        this.snackBar.open(`"${this.tour?.name}" added to cart!`, 'View Cart', { duration: 4000 })
          .onAction().subscribe(() => this.router.navigate(['/shopping-cart']));
        
        this.cartStateService.updateCartCount(updatedCart.items.length);
        this.hasPurchased = true; // Status se odmah ažurira nakon uspešne kupovine
        this.isAddingToCart = false;
      },
      error: (err) => {
        const errorMessage = err.error?.message || 'Failed to add item to cart.';
        this.snackBar.open(errorMessage, 'Dismiss', { duration: 5000 });
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
