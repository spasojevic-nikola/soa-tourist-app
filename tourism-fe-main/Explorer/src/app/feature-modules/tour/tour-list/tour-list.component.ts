import { Component, OnInit } from '@angular/core';
import { Observable } from 'rxjs';
import { TourService } from '../tour.service';
import { Tour } from '../model/tour.model';
import { Router } from '@angular/router';
import { KeypointService } from '../../tour-keypoints/keypoint.service';
import { KeyPoint } from '../../tour-keypoints/model/keypoint.model';

@Component({
  selector: 'xp-tour-list',
  templateUrl: './tour-list.component.html',
  styleUrls: ['./tour-list.component.css']
})
export class TourListComponent implements OnInit {
  tours$: Observable<Tour[]>;
  keypointsMap: Map<number, boolean> = new Map(); 

  constructor(private tourService: TourService, 
    private keypointService: KeypointService,
    private router: Router) { }

  ngOnInit(): void {
    this.tours$ = this.tourService.getAuthorTours();
    this.loadKeypointsInfo();
  }

  loadKeypointsInfo(): void {
    this.tours$.subscribe(tours => {
      tours.forEach(tour => {
        this.keypointService.getKeyPointsByTour(tour.id).subscribe(keypoints => {
          this.keypointsMap.set(tour.id, keypoints.length > 0);
        });
      });
    });
  }

  hasKeypoints(tourId: number): boolean {
    return this.keypointsMap.get(tourId) || false;
  }

  navigateToKeyPoints(tourId: number): void {
    this.router.navigate(['/keypoints', tourId]);
  }
}