import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { TourService } from '../tour.service';
import { CreateTourPayload } from '../dto/tour-creation.dto';
import { CreateKeyPointPayload } from '../../tour-keypoints/model/keypoint.model';

@Component({
  selector: 'xp-tour-wizard',
  templateUrl: './tour-wizard.component.html',
  styleUrls: ['./tour-wizard.component.css']
})
export class TourWizardComponent {
  currentStep = 1;
  tourData: CreateTourPayload | null = null;
  keyPoints: CreateKeyPointPayload[] = [];
  isSubmitting = false;

  constructor(
    private tourService: TourService,
    private router: Router
  ) {}

  onStep1Completed(tourData: CreateTourPayload): void {
    this.tourData = tourData;
    this.currentStep = 2;
  }

  onStep2Completed(keyPoints: CreateKeyPointPayload[]): void {
    this.keyPoints = keyPoints;
    this.createTour();
  }

  createTour(): void {
    if (!this.tourData || this.keyPoints.length === 0) {
      alert('Please add at least one key point before creating the tour.');
      return;
    }

    this.isSubmitting = true;

    // Kreiraj payload koji odgovara backend DTO
    const payload = {
      name: this.tourData.name,
      description: this.tourData.description,
      difficulty: this.tourData.difficulty,
      tags: this.tourData.tags,
      keyPoints: this.keyPoints
    };

    this.tourService.createTour(payload).subscribe({
      next: (createdTour) => {
        this.isSubmitting = false;
        alert('Tour created successfully in Draft status!');
        this.router.navigate(['/tours']); // Navigate to tour list
      },
      error: (err) => {
        this.isSubmitting = false;
        console.error('Error creating tour:', err);
        alert('Failed to create tour. Please try again.');
      }
    });
  }

  goBack(): void {
    if (this.currentStep > 1) {
      this.currentStep--;
    } else {
      this.router.navigate(['/tours']);
    }
  }
}