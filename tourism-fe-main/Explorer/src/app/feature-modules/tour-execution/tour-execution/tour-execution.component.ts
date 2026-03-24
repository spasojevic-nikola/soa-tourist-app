import { Component, OnInit, OnDestroy } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { TourExecutionService } from '../tour-execution.service';
import { TourService } from '../../tour/tour.service';
import { PositionSimulatorService } from '../../position-simulator/position-simulator.service';

@Component({
  selector: 'xp-tour-execution',
  templateUrl: './tour-execution.component.html',
  styleUrls: ['./tour-execution.component.css']
})
export class TourExecutionComponent implements OnInit, OnDestroy {
  executionId: number;
  execution: any;
  tour: any;
  positionInterval: any;
  currentPosition: any;
  isLoading: boolean = true;
  
  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private tourExecutionService: TourExecutionService,
    private tourService: TourService,
    private positionService: PositionSimulatorService
  ) {
    this.executionId = +this.route.snapshot.params['id'];
  }

  ngOnInit(): void {
    this.loadExecution();
    this.startPositionTracking();
    this.loadCurrentPosition();
  }

  ngOnDestroy(): void {
    if (this.positionInterval) {
      clearInterval(this.positionInterval);
    }
  }

  loadExecution(): void {
    this.tourExecutionService.getExecutionDetails(this.executionId).subscribe({
      next: (execution) => {
        this.execution = execution;
        this.loadTourDetails(execution.tourId);
      },
      error: (err) => {
        console.error('Error loading execution:', err);
        this.isLoading = false;
      }
    });
  }

  loadTourDetails(tourId: number): void {
    this.tourService.getTourById(tourId).subscribe({
      next: (tour) => {
        this.tour = tour;
        this.isLoading = false;
        console.log('Tour loaded:', tour);
      },
      error: (err) => {
        console.error('Error loading tour:', err);
        this.isLoading = false;
      }
    });
  }

  loadCurrentPosition(): void {
    this.currentPosition = this.positionService.getCurrentPosition();
    console.log('Current position loaded:', this.currentPosition);
  }

  startPositionTracking(): void {
    // Prvo odmah proveri poziciju
    this.checkPosition();
    
    // Onda šalji na svakih 10 sekundi
    this.positionInterval = setInterval(() => {
      this.checkPosition();
    }, 10000);
  }

checkPosition(): void {
  this.loadCurrentPosition();
  
  if (!this.currentPosition) {
    console.log(' No current position available');
    return;
  }

  console.log('🔍 Checking position:', {
    executionId: this.executionId,
    position: this.currentPosition
  });

  this.tourExecutionService.checkPosition(
    this.executionId, 
    this.currentPosition.lat, 
    this.currentPosition.lng
  ).subscribe({
    next: (completedKeyPoints: number[]) => { // DIREKTNO ARRAY, NE OBJEKAT
      console.log('✅ Position check successful. Completed key points:', completedKeyPoints);
      
      // Ažuriraj execution sa novim completed key points
      if (completedKeyPoints && completedKeyPoints.length > 0) {
        console.log('🎉 New completed key points:', completedKeyPoints);
        
        // Ažuriraj listu completed key points (spreči duplikate)
        const newPoints = completedKeyPoints.filter((kpId: number) => 
          !this.execution.completedKeyPoints.includes(kpId)
        );
        
        if (newPoints.length > 0) {
          this.execution.completedKeyPoints = [...this.execution.completedKeyPoints, ...newPoints];
          console.log('📊 Updated completed key points:', this.execution.completedKeyPoints);
          
          // Osveži prikaz
          this.refreshExecution();
        }
      }
    },
    error: (error) => {
      console.error('❌ Error checking position:', error);
      console.log('🔍 Full error details:', {
        status: error.status,
        statusText: error.statusText,
        error: error.error,
        url: error.url
      });
      
      if (error.status === 404) {
        alert('Tour execution not found. Redirecting to tours...');
        this.router.navigate(['/tours']);
      }
    }
  });
}

canCompleteTour(): boolean {
  if (!this.execution || !this.tour || !this.tour.keyPoints) return false;
  
  // Ne može da kompletira ako je već kompletirana
  if (this.execution.status === 'COMPLETED') return false;
  
  // Može da kompletira samo ako su sve tačke završene
  const totalKeyPoints = this.tour.keyPoints.length;
  const completedKeyPoints = this.execution.completedKeyPoints?.length || 0;
  
  return completedKeyPoints >= totalKeyPoints;
}

  completeTour(): void {
    if (confirm('Are you sure you want to complete this tour?')) {
      this.tourExecutionService.completeTour(this.executionId).subscribe({
        next: (response) => {
          console.log('✅ Tour completed:', response);
          alert('Tour completed successfully! 🎉');
          this.execution.status = 'COMPLETED';
          this.router.navigate(['/tours', this.tour.id]);
        },
        error: (err) => {
          console.error('Error completing tour:', err);
          alert('Error completing tour. Please try again.');
        }
      });
    }
  }

  abandonTour(): void {
    if (confirm('Are you sure you want to abandon this tour?')) {
      this.tourExecutionService.abandonTour(this.executionId).subscribe({
        next: (response) => {
          console.log('✅ Tour abandoned:', response);
          alert('Tour abandoned');
          this.execution.status = 'ABANDONED';
          this.router.navigate(['/tours', this.tour.id]);
        },
        error: (err) => {
          console.error('Error abandoning tour:', err);
          alert('Error abandoning tour. Please try again.');
        }
      });
    }
  }

  getProgressPercentage(): number {
    if (!this.tour || !this.tour.keyPoints || this.tour.keyPoints.length === 0) return 0;
    const total = this.tour.keyPoints.length;
    const completed = this.execution?.completedKeyPoints?.length || 0;
    const percentage = Math.round((completed / total) * 100);
    console.log(`📊 Progress: ${completed}/${total} = ${percentage}%`);
    return percentage;
  }

  isKeyPointCompleted(keyPointId: number): boolean {
    const isCompleted = this.execution?.completedKeyPoints?.includes(keyPointId) || false;
    console.log(`🔍 KeyPoint ${keyPointId} completed:`, isCompleted);
    return isCompleted;
  }

  // Dodatna metoda za debug
  refreshExecution(): void {
    this.loadExecution();
    console.log('🔄 Execution refreshed');
  }
}