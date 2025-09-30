import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatIconModule } from '@angular/material/icon';
import { ReactiveFormsModule } from '@angular/forms';
import { TourCreateComponent } from './create-tour/tour-create/tour-create.component';
import { TourRoutingModule } from './tour-routing.module';
import { TourListComponent } from './tour-list/tour-list.component';
import { TourWizardComponent } from './tour-wizard/tour-wizard.component';
import { TourKeypointsModule } from '../tour-keypoints/tour-keypoints.module';

@NgModule({
  declarations: [
    TourCreateComponent,
    TourListComponent,
    TourWizardComponent
  ],
  imports: [
    CommonModule,
    MatIconModule,
    ReactiveFormsModule, 
    TourRoutingModule,  
    TourKeypointsModule  
  ]
})
export class TourModule { }
