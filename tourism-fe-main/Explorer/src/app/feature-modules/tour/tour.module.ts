import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatIconModule } from '@angular/material/icon';
import { MatDialogModule } from '@angular/material/dialog';
import { MatButtonModule } from '@angular/material/button';
import { ReactiveFormsModule } from '@angular/forms';
import { TourCreateComponent } from './create-tour/tour-create/tour-create.component';
import { TourRoutingModule } from './tour-routing.module';
import { TourListComponent } from './tour-list/tour-list.component';
import { TourWizardComponent } from './tour-wizard/tour-wizard.component';
import { TourKeypointsModule } from '../tour-keypoints/tour-keypoints.module';
import { TourDetailsComponent } from './tour-details/tour-details.component';
import { ReviewDialogComponent } from './review-dialog/review-dialog.component';
import { MatSnackBarModule } from '@angular/material/snack-bar'; 
import { HttpClientModule } from '@angular/common/http'; 
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner'; 


@NgModule({
  declarations: [
    TourCreateComponent,
    TourListComponent,
    TourWizardComponent,
    TourDetailsComponent,
    ReviewDialogComponent
  ],
  imports: [
    CommonModule,
    MatIconModule,
    MatDialogModule,
    MatButtonModule,
    ReactiveFormsModule, 
    TourRoutingModule,  
    TourKeypointsModule,
    MatSnackBarModule,
    HttpClientModule,
    MatProgressSpinnerModule

  ]
})
export class TourModule { }
