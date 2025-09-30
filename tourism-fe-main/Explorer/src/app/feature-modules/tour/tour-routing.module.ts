import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { TourListComponent } from './tour-list/tour-list.component';
import { TourWizardComponent  } from './tour-wizard/tour-wizard.component'

const routes: Routes = [
  {
    path: 'create',
    component: TourWizardComponent  
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