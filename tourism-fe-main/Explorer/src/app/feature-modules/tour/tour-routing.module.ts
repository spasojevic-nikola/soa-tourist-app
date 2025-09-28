import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { TourCreateComponent } from './create-tour/tour-create/tour-create.component';
import { TourListComponent } from './tour-list/tour-list.component';

const routes: Routes = [
  {
    path: 'create',
    component: TourCreateComponent
  },
  {
    path: '',
    component: TourListComponent
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class TourRoutingModule { }