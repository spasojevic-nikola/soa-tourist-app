import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { BlogCreationComponent } from './blog-creation/blog-creation.component';
import { MatIconModule } from '@angular/material/icon';
import { ReactiveFormsModule } from '@angular/forms';
import { BlogRoutingModule } from './blog-routing.module';



@NgModule({
  declarations: [
    BlogCreationComponent
  ],
  imports: [
    CommonModule,
    MatIconModule,
    ReactiveFormsModule,
    BlogRoutingModule 
  ]
})
export class BlogModule { }
