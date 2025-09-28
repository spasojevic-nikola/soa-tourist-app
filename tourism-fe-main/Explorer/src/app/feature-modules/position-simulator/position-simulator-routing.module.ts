import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { PositionSimulatorComponent } from './position-simulator.component';

const routes: Routes = [
  { path: '', component: PositionSimulatorComponent }
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class PositionSimulatorRoutingModule { }