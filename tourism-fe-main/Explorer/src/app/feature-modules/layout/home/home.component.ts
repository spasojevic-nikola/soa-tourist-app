import { Component, OnInit } from '@angular/core';
import { TourService } from '../../tour/tour.service';
import { Tour } from '../../tour/model/tour.model';

@Component({
  selector: 'xp-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.css']
})
export class HomeComponent implements OnInit {
  publishedTours: Tour[] = [];
  isLoading = true;

  constructor(private tourService: TourService) {}

  ngOnInit(): void {
    this.loadPublishedTours();
  }

  loadPublishedTours(): void {
    this.isLoading = true;
    this.tourService.getAllPublishedTours().subscribe({
      next: (tours) => {
        this.publishedTours = tours;
        this.isLoading = false;
      },
      error: (err) => {
        console.error('Error loading published tours:', err);
        this.isLoading = false;
      }
    });
  }
}
