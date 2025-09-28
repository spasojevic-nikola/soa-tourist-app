import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatIconModule } from '@angular/material/icon';
import { ReactiveFormsModule } from '@angular/forms';
import { TourCreateComponent } from './create-tour/tour-create/tour-create.component';
import { TourRoutingModule } from './tour-routing.module';
import { TourListComponent } from './tour-list/tour-list.component';



@NgModule({
  declarations: [
    TourCreateComponent,
    TourListComponent
  ],
  imports: [
    CommonModule,
    MatIconModule,
    ReactiveFormsModule, 
    TourRoutingModule,    
  ]
})
export class TourModule { }
