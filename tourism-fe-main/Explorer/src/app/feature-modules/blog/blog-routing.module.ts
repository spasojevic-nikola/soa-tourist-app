import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { BlogCreationComponent } from './blog-creation/blog-creation.component';

const routes: Routes = [
  {
    path: 'create',
    component: BlogCreationComponent
  }
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class BlogRoutingModule { }