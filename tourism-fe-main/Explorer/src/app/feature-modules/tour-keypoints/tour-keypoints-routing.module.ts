import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { TourKeypointsComponent } from './tour-keypoints/tour-keypoints.component';
import { TourMapCreationComponent } from './tour-map-creation/tour-map-creation.component';

const routes: Routes = [
  {
    path: ':tourId',
    component: TourKeypointsComponent
  },
  {
    path: ':tourId/map',       
    component: TourMapCreationComponent
  }
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class KeypointsRoutingModule { }