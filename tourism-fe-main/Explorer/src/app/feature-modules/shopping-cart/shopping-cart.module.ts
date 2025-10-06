import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';

import { ShoppingCartRoutingModule } from './shopping-cart-routing.module'; 
import { CartComponent } from './components/cart/cart.component'; 

import { MatProgressSpinnerModule } from '@angular/material/progress-spinner'; 
import { MatIconModule } from '@angular/material/icon';
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { MatDividerModule } from '@angular/material/divider'; 
import { MatSnackBarModule } from '@angular/material/snack-bar'; 



@NgModule({
  declarations: [
    CartComponent 
  ],
  imports: [
    CommonModule,
    ShoppingCartRoutingModule,
    MatProgressSpinnerModule,
    MatSnackBarModule,
    MatIconModule,
    MatCardModule,
    MatButtonModule,
    MatDividerModule
  ]
})
export class ShoppingCartModule { } 