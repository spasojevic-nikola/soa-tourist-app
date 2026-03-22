import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatIconModule } from '@angular/material/icon';
import { MatButtonModule } from '@angular/material/button';
import { MatDialogModule } from '@angular/material/dialog';
import { MatSnackBarModule } from '@angular/material/snack-bar';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatMenuModule } from '@angular/material/menu';
import { ReactiveFormsModule } from '@angular/forms';
import { TourCreateComponent } from './create-tour/tour-create/tour-create.component';
import { TourRoutingModule } from './tour-routing.module';
import { TourListComponent } from './tour-list/tour-list.component';
import { TourWizardComponent } from './tour-wizard/tour-wizard.component';
import { TourKeypointsModule } from '../tour-keypoints/tour-keypoints.module';
import { TourDetailsComponent } from './tour-details/tour-details.component';
import { ReviewDialogComponent } from './review-dialog/review-dialog.component';
import { EditTourDialogComponent } from './edit-tour-dialog/edit-tour-dialog.component';

@NgModule({
  declarations: [
    TourCreateComponent,
    TourListComponent,
    TourWizardComponent,
    TourDetailsComponent,
    ReviewDialogComponent,
    EditTourDialogComponent
  ],
  imports: [
    CommonModule,
    MatIconModule,
    MatButtonModule,
    MatDialogModule,
    MatSnackBarModule,
    MatProgressSpinnerModule,
    MatMenuModule,
    ReactiveFormsModule, 
    TourRoutingModule,  
    TourKeypointsModule  
  ]
})
export class TourModule { }
