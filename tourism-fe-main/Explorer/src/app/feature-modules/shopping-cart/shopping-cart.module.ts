import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';

import { ShoppingCartRoutingModule } from './shopping-cart-routing.module'; 

import { CartComponent } from './components/cart/cart.component'; 


@NgModule({
  declarations: [
    CartComponent 
  ],
  imports: [
    CommonModule,
    ShoppingCartRoutingModule 
  ]
})
export class ShoppingCartModule { } 