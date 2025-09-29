import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { TourKeypointsComponent } from './tour-keypoints/tour-keypoints.component';

const routes: Routes = [
  {
    path: ':tourId',
    component: TourKeypointsComponent
  }
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class KeypointsRoutingModule { }