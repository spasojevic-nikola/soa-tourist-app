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
import { AuthService } from '../../../infrastructure/auth/auth.service';

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

  isTourPurchased: boolean = false;
  checkingPurchaseStatus: boolean = true;


    // Tour Execution polja
    hasActiveExecution: boolean = false;
    activeExecutionId: number | null = null;
    checkingExecutionStatus: boolean = false;
    executionProgress: number = 0; 
    hasCompletedExecution: boolean = false;

  constructor(
    private route: ActivatedRoute,
    private tourService: TourService,
    private reviewService: ReviewService,
    private mapService: MapService,
    private dialog: MatDialog,
    private router: Router,
    private cartService: CartService, 
    private cartStateService: CartStateService,
    private snackBar: MatSnackBar,
    private authService: AuthService
  ) {}

  ngOnInit(): void {
    const tourId = this.route.snapshot.params['id'];
    this.loadTourDetails(tourId);
    this.loadReviews(tourId);
    this.loadReviewStats(tourId);
    this.checkAllExecutions(tourId);
    this.checkPurchaseStatus(tourId); 

  }
  checkPurchaseStatus(tourId: number): void {
    this.checkingPurchaseStatus = true;
    this.cartService.hasPurchased(String(tourId)).subscribe({
      next: (response) => {
        this.isTourPurchased = response.isPurchased;
        this.checkingPurchaseStatus = false;
      },
      error: (err) => {
        console.error('Error checking purchase status:', err);
        // U slučaju greške, pretpostavljamo da nije kupljeno da ne bi prikazali pogrešne opcije
        this.isTourPurchased = false; 
        this.checkingPurchaseStatus = false;
      }
    });
  }


  loadTourDetails(tourId: number): void {
  this.tourService.getTourById(tourId).subscribe({
    next: async (tour) => {
      this.tour = tour;
      
      // Učitaj adrese za sve keypoints
      if (tour.keyPoints && tour.keyPoints.length > 0) {
        for (let keypoint of tour.keyPoints) {
          if (!keypoint.address) {
            keypoint.address = await this.mapService.reverseGeocode(
              keypoint.latitude,
              keypoint.longitude
            );
          }
        }
        
        // Postavi adresu prvog keypointa za prikaz kada nije kupljeno
        this.keypointAddress = tour.keyPoints[0].address || 'Loading address...';
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

    if (!this.tour || !this.tour.id) { 
        this.snackBar.open('Tour details are incomplete or still loading.', 'Dismiss', { duration: 3000 });
        return;
    }
    
    if (!isPublished ) {
        this.snackBar.open('Cannot purchase: The tour must be published.', 'Dismiss', { duration: 4000 });
        return;
    }
  
    this.isAddingToCart = true; 
  
    // Kreiramo objekat koji sadrži SAMO ID ture
    const itemToAdd = {
        tourId: String(this.tour.id)
    };
  
    // 3. POZIV BACKENDA (sada sa ispravnim, "glupim" objektom)
    this.cartService.addItem(itemToAdd).subscribe({
        next: (updatedCart) => {
            // Pošto 'itemToAdd' više nema ime, koristimo 'this.tour.name' za poruku
            this.snackBar.open(`"${this.tour!.name}" added to cart!`, 'View Cart', { duration: 4000 })
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

  isAuthor(): boolean {
    const currentUser = this.authService.user$.getValue();
    return this.tour?.authorId === currentUser.id;
  }

   checkAllExecutions(tourId: number): void {
  this.checkingExecutionStatus = true;
  
  this.tourService.getAllExecutionsForTour(tourId).subscribe({
    next: (executions) => {
      console.log('🔍 All executions:', executions);
      
      const activeExecution = executions.find(e => e.status === 'STARTED');
      const completedExecution = executions.find(e => e.status === 'COMPLETED');
      
      this.hasActiveExecution = !!activeExecution;
      this.hasCompletedExecution = !!completedExecution;
      this.activeExecutionId = activeExecution?.id || completedExecution?.id || null;
      
      if (activeExecution && this.tour?.keyPoints) {
        const totalKeyPoints = this.tour.keyPoints.length;
        const completed = activeExecution.completedKeyPoints?.length || 0;
        this.executionProgress = Math.round((completed / totalKeyPoints) * 100);
      }
      
      this.checkingExecutionStatus = false;
    },
    error: (err) => {
      console.error('Error checking executions:', err);
      this.hasActiveExecution = false;
      this.hasCompletedExecution = false;
      this.checkingExecutionStatus = false;
      this.executionProgress = 0;  
    }
  });
}

  startTourExecution(): void {
    if (!this.tour) return;

    // Prvo dohvati trenutnu lokaciju iz Position Simulatora
    this.tourService.startTourExecution(this.tour.id).subscribe({
      next: (execution) => {
        this.activeExecutionId = execution.id;
        this.hasActiveExecution = true;
        
        // Redirect na tour execution stranicu
        this.router.navigate(['/tour-execution', execution.id]);
      },
      error: (err) => {
        console.error('Error starting tour:', err);
        console.error('Full error:', err.error); 
        alert('Failed to start tour: ' + err.message);
      }
    });
  }
}

