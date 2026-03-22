import { Component, OnInit } from '@angular/core';
import { Tour } from '../../tour/model/tour.model';
import { TourService } from '../../tour/tour.service';

@Component({
  selector: 'xp-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.css']
})
export class HomeComponent implements OnInit {
  publishedTours: Tour[] = [];
  isLoadingPublished = false;
  publishedToursError = '';

  constructor(private tourService: TourService) {}

  ngOnInit(): void {
    this.loadPublishedTours();
  }

  private loadPublishedTours(): void {
    this.isLoadingPublished = true;
    this.publishedToursError = '';

    this.tourService.getPublishedTours().subscribe({
      next: (tours: Tour[]) => {
        this.publishedTours = tours;
        this.isLoadingPublished = false;
      },
      error: (err: unknown) => {
        console.error('Failed to load published tours', err);
        this.publishedTours = [];
        this.publishedToursError = 'Failed to load published tours. Please try again.';
        this.isLoadingPublished = false;
      }
    });
  }

  getTourSummary(tour: Tour): string {
    const description = (tour.description || '').trim();

    if (!description) {
      return 'Author has not added a summary yet.';
    }

    const maxLength = 220;
    if (description.length <= maxLength) {
      return description;
    }

    return description.substring(0, maxLength).trimEnd() + '…';
  }
}
