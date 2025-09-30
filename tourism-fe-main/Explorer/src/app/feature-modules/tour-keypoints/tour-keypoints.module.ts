import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { KeypointsRoutingModule } from './tour-keypoints-routing.module'
import { TourKeypointsComponent } from './tour-keypoints/tour-keypoints.component'
import { ReactiveFormsModule } from '@angular/forms';
import { TourMapCreationComponent } from './tour-map-creation/tour-map-creation.component';

@NgModule({
  declarations: [
    TourKeypointsComponent,
    TourMapCreationComponent
  ],
  imports: [
    CommonModule,
    KeypointsRoutingModule,
    ReactiveFormsModule
  ],
  exports: [
    TourMapCreationComponent   
  ]
})
export class TourKeypointsModule { }
