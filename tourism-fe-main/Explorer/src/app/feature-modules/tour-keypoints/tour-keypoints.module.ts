import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { KeypointsRoutingModule } from './tour-keypoints-routing.module'
import { TourKeypointsComponent } from './tour-keypoints/tour-keypoints.component'
import { ReactiveFormsModule } from '@angular/forms';

@NgModule({
  declarations: [
    TourKeypointsComponent
  ],
  imports: [
    CommonModule,
    KeypointsRoutingModule,
    ReactiveFormsModule
  ]
})
export class TourKeypointsModule { }
