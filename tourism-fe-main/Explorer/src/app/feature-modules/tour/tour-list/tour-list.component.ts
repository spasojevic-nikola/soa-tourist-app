import { Component, OnInit } from '@angular/core';
import { Observable } from 'rxjs';
import { TourService } from '../tour.service';
import { Tour } from '../model/tour.model';

@Component({
  selector: 'xp-tour-list',
  templateUrl: './tour-list.component.html',
  styleUrls: ['./tour-list.component.css']
})
export class TourListComponent implements OnInit {
  tours$: Observable<Tour[]>;

  constructor(private tourService: TourService) { }

  ngOnInit(): void {
    this.tours$ = this.tourService.getAuthorTours();
  }
}