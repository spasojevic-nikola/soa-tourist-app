import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { PositionSimulatorComponent } from './position-simulator.component';
import { PositionSimulatorRoutingModule } from './position-simulator-routing.module';
import { MaterialModule } from 'src/app/infrastructure/material/material.module';
import { FormsModule } from '@angular/forms';

@NgModule({
  declarations: [
    PositionSimulatorComponent
  ],
  imports: [
    CommonModule,
    PositionSimulatorRoutingModule,
    MaterialModule,
    FormsModule
  ]
})
export class PositionSimulatorModule { }