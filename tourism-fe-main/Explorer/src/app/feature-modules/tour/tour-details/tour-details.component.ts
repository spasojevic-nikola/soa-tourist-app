import { AfterViewInit, Component, ElementRef, OnDestroy, OnInit, ViewChild } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { MatDialog } from '@angular/material/dialog';
import { TourService } from '../tour.service';
import { ReviewService } from '../review.service';
import { Tour } from '../model/tour.model';
import { Review, ReviewStats } from '../model/review.model';
import { MapService } from '../services/map-service.service';
import { ReviewDialogComponent } from '../review-dialog/review-dialog.component';
import { EditTourDialogComponent } from '../edit-tour-dialog/edit-tour-dialog.component';
import { UpdateTourPayload } from '../dto/tour-creation.dto';
import { KeypointDialogService } from '../services/keypoint-dialog.service';
import { KeypointService } from '../../tour-keypoints/keypoint.service';
import { CreateKeyPointPayload, KeyPoint } from '../../tour-keypoints/model/keypoint.model';
import { CartService } from '../../shopping-cart/services/cart.service';
import { MatSnackBar } from '@angular/material/snack-bar';
import { CartStateService } from '../../shopping-cart/services/cart-state.service';
import { AuthService } from '../../../infrastructure/auth/auth.service';
import { StakeholdersService } from 'src/app/infrastructure/stakeholders.service';
import { User as StakeholderUser } from '../../user-profile/profile/model/profile.model';
import { Subscription } from 'rxjs';

@Component({
  selector: 'xp-tour-details',
  templateUrl: './tour-details.component.html',
  styleUrls: ['./tour-details.component.css']
})
export class TourDetailsComponent implements OnInit, OnDestroy, AfterViewInit {
  tour: Tour | null = null;
  isLoading = true;
  keypointAddress: string = '';
  reviews: Review[] = [];
  reviewStats: ReviewStats | null = null;
  loadingReviews = false;

  isAddingToCart: boolean = false; 

  isTourPurchased: boolean = false;
  checkingPurchaseStatus: boolean = true;

  @ViewChild('keypointCarousel') keypointCarousel?: ElementRef<HTMLDivElement>;
  canScrollPrev = false;
  canScrollNext = false;

  authorInfo: StakeholderUser | null = null;
  authorInfoLoading = true;
  authorInfoError = false;

  private authSubscription?: Subscription;
  private authorSubscription?: Subscription;
  private currentUserRole = '';
  private currentUserId = 0;


    // Tour Execution polja
    hasActiveExecution: boolean = false;
    activeExecutionId: number | null = null;
    checkingExecutionStatus: boolean = false;
    executionProgress: number = 0; 
    hasCompletedExecution: boolean = false;

    private activeExecutionData: any = null;

    private calculateExecutionProgress(): void {
        if (this.activeExecutionData && this.tour?.keyPoints) {
            const totalKeyPoints = this.tour.keyPoints.length;
            const completed = this.activeExecutionData.completedKeyPoints?.length || 0;
            this.executionProgress = Math.round((completed / totalKeyPoints) * 100);
        }
    }

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
    private authService: AuthService,
    private stakeholdersService: StakeholdersService,
    private keypointDialogService: KeypointDialogService,
    private keypointService: KeypointService
  ) {
    this.authSubscription = this.authService.user$.subscribe(user => {
      this.currentUserRole = (user.role || '').toLowerCase();
      this.currentUserId = user.id;
    });
  }

  ngOnInit(): void {
    const tourIdParam = this.route.snapshot.paramMap.get('id');
    const tourId = tourIdParam ? Number(tourIdParam) : NaN;

    if (Number.isNaN(tourId)) {
      this.router.navigate(['/tours']);
      return;
    }

    this.loadTourDetails(tourId);
    this.loadReviews(tourId);
    this.loadReviewStats(tourId);
    this.checkAllExecutions(tourId);
    this.checkPurchaseStatus(tourId); 

  }

  ngAfterViewInit(): void {
    this.updateCarouselNav();
  }

  ngOnDestroy(): void {
    this.authSubscription?.unsubscribe();
    this.authorSubscription?.unsubscribe();
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
        // U sluÄaju greÅ¡ke, pretpostavljamo da nije kupljeno da ne bi prikazali pogreÅ¡ne opcije
        this.isTourPurchased = false; 
        this.checkingPurchaseStatus = false;
      }
    });
  }


  loadTourDetails(tourId: number): void {
    this.tourService.getTourById(tourId).subscribe({
      next: async (tour) => {
        this.tour = tour;
        this.loadAuthorInfo(tour.authorId);
        
        if (tour.keyPoints && tour.keyPoints.length > 0) {
          for (let keypoint of tour.keyPoints) {
            if (!keypoint.address) {
              keypoint.address = await this.mapService.reverseGeocode(
                keypoint.latitude,
                keypoint.longitude
              );
            }
          }
          this.keypointAddress = tour.keyPoints[0].address || 'Loading address...';
        }
        
        this.calculateExecutionProgress();
        this.isLoading = false;
        setTimeout(() => this.updateCarouselNav(), 0);
      },
      error: (err) => {
        console.error('Error loading tour details:', err);
        this.isLoading = false;
      }
    });
  }

  private loadAuthorInfo(authorId: number): void {
    if (!authorId) {
      this.authorInfoLoading = false;
      return;
    }

    this.authorInfoLoading = true;
    this.authorInfoError = false;
    this.authorSubscription?.unsubscribe();
    this.authorSubscription = this.stakeholdersService.getUserById(authorId).subscribe({
      next: (user) => {
        this.authorInfo = user;
        this.authorInfoLoading = false;
      },
      error: (err) => {
        console.error('Error loading author info:', err);
        this.authorInfo = null;
        this.authorInfoError = true;
        this.authorInfoLoading = false;
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

  onPublishTour(): void {
    if (!this.tour || !this.tour.id) return;
    
    this.tourService.publishTour(this.tour.id).subscribe({
      next: (updatedTour) => {
        this.tour = updatedTour;
        this.snackBar.open('Tour published successfully!', 'Close', { duration: 3000 });
      },
      error: (err) => {
        console.error('Error publishing tour:', err);
        this.snackBar.open(err.error?.message || 'Failed to publish tour.', 'Close', { duration: 5000 });
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
  
    // Kreiramo objekat koji sadrÅ¾i SAMO ID ture
    const itemToAdd = {
        tourId: String(this.tour.id)
    };
  
    // 3. POZIV BACKENDA (sada sa ispravnim, "glupim" objektom)
    this.cartService.addItem(itemToAdd).subscribe({
        next: (updatedCart) => {
            // PoÅ¡to 'itemToAdd' viÅ¡e nema ime, koristimo 'this.tour.name' za poruku
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

  onEditKeyPoint(keypoint: any): void {
    const mapToPayload = (kp: any): CreateKeyPointPayload => ({
      name: kp.name,
      description: kp.description,
      latitude: kp.latitude,
      longitude: kp.longitude,
      image: (kp.image as unknown as File | null) || null,
      order: kp.order,
      address: kp.address
    });

    this.keypointDialogService.openKeypointDialog(
      keypoint.latitude,
      keypoint.longitude,
      keypoint.order,
      mapToPayload(keypoint)
    ).subscribe((result) => {
      if (result) {
        this.keypointService.updateKeyPoint(keypoint.id, result).subscribe({
          next: () => {
             this.snackBar.open('Key point updated successfully', 'Close', { duration: 3000 });
             this.loadTourDetails(this.tour!.id);
          },
          error: (err) => {
             console.error('Error updating key point:', err);
             this.snackBar.open('Failed to update key point', 'Close', { duration: 3000 });
          }
        });
      }
    });
  }

  onDeleteKeyPoint(keypointId: number): void {
    if (confirm('Are you sure you want to delete this key point?')) {
      this.keypointService.deleteKeyPoint(keypointId).subscribe({
        next: () => {
          this.snackBar.open('Key point deleted successfully', 'Close', { duration: 3000 });
          this.loadTourDetails(this.tour!.id);
        },
        error: (err) => {
          console.error('Error deleting key point:', err);
          this.snackBar.open('Failed to delete key point', 'Close', { duration: 3000 });
        }
      });
    }
  }

  onEditTour(): void {
    if (!this.tour) return;

    const dialogRef = this.dialog.open(EditTourDialogComponent, {
      width: '600px',
      data: { tour: this.tour }
    });

    dialogRef.afterClosed().subscribe((result: UpdateTourPayload) => {
      if (result) {
        this.tourService.updateTour(this.tour!.id, result).subscribe({
          next: () => {
             this.snackBar.open('Tour updated successfully', 'Close', { duration: 3000 });
             this.loadTourDetails(this.tour!.id);
          },
          error: (err) => {
             console.error('Error updating tour:', err);
             this.snackBar.open('Failed to update tour', 'Close', { duration: 3000 });
          }
        });
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

  onEditReview(review: Review): void {
    if (!this.tour) return;

    const dialogRef = this.dialog.open(ReviewDialogComponent, {
      width: '600px',
      data: {
        tourId: this.tour.id,
        tourName: this.tour.name,
        review: review
      }
    });

    dialogRef.afterClosed().subscribe(result => {
      if (result) {
        this.reviewService.updateReview(review.id, result).subscribe({
          next: () => {
            this.loadReviews(this.tour!.id);
            this.loadReviewStats(this.tour!.id);
            this.snackBar.open('Review updated successfully.', 'Close', { duration: 3000 });
          },
          error: (err) => {
            console.error('Error updating review:', err);
            this.snackBar.open('Failed to update review.', 'Close', { duration: 3000 });
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
    return this.tour?.authorId === this.currentUserId;
  }

  isReviewAuthor(review: Review): boolean {
    return review.touristId === this.currentUserId;
  }

  canViewFullTour(): boolean {
    if (!this.tour) {
      return false;
    }

    if (this.isAuthor()) {
      return true;
    }

    const allowedRoles = ['guide', 'tourist'];
    return this.tour.status === 'Published' && allowedRoles.includes(this.currentUserRole);
  }

   checkAllExecutions(tourId: number): void {
  this.checkingExecutionStatus = true;
  
  this.tourService.getAllExecutionsForTour(tourId).subscribe({
    next: (executions) => {
      console.log('ðŸ” All executions:', executions);
      
      const activeExecution = executions.find(e => e.status === 'STARTED');
      const completedExecution = executions.find(e => e.status === 'COMPLETED');
      
      this.hasActiveExecution = !!activeExecution;
        this.activeExecutionData = activeExecution;
      this.hasCompletedExecution = !!completedExecution;
      this.activeExecutionId = activeExecution?.id || completedExecution?.id || null;
      
      this.calculateExecutionProgress();
      
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

  scrollKeypoints(direction: 'prev' | 'next'): void {
    const container = this.keypointCarousel?.nativeElement;
    if (!container) return;

    const scrollAmount = container.clientWidth * 0.85;
    container.scrollBy({
      left: direction === 'next' ? scrollAmount : -scrollAmount,
      behavior: 'smooth'
    });

    setTimeout(() => this.updateCarouselNav(), 350);
  }

  onCarouselScroll(): void {
    this.updateCarouselNav();
  }

  private updateCarouselNav(): void {
    const container = this.keypointCarousel?.nativeElement;

    if (!container) {
      this.canScrollPrev = false;
      this.canScrollNext = false;
      return;
    }

    const { scrollLeft, scrollWidth, clientWidth } = container;
    this.canScrollPrev = scrollLeft > 8;
    this.canScrollNext = scrollLeft + clientWidth < scrollWidth - 8;
  }
}




