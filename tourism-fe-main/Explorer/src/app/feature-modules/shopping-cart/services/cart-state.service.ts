import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable } from 'rxjs';
import { CartService } from './cart.service';

@Injectable({
  providedIn: 'root'
})
export class CartStateService {

  private cartItemCountSubject = new BehaviorSubject<number>(0);
  
  // 💡 Observable koji komponente mogu da "prate" (subscribe)
  public cartItemCount$: Observable<number> = this.cartItemCountSubject.asObservable();

  constructor(private cartService: CartService) { 
    // 💡 Učitaj inicijalnu vrednost iz backenda
    this.initializeCartCount();
  }
  private initializeCartCount(): void {
    this.cartService.getCart().subscribe({
      next: (cart) => {
        this.cartItemCountSubject.next(cart.items.length);
      },
      error: (err) => {
        console.error('Failed to initialize cart count:', err);
        // Ostavi na 0 ako korisnik nije ulogovan ili ima problem sa backendom
        this.cartItemCountSubject.next(0);
      }
    });
  }

  updateCartCount(newCount: number): void {
    this.cartItemCountSubject.next(newCount);
  }

  refreshCartCount(): void {
    this.initializeCartCount();
  }
}
