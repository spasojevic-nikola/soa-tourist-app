import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class CartStateService {

  private cartItemCountSubject = new BehaviorSubject<number>(0);
  
  // 💡 Observable koji komponente mogu da "prate" (subscribe)
  public cartItemCount$: Observable<number> = this.cartItemCountSubject.asObservable();

  constructor() { 
    // Opcionalno: Možete ovde inicijalno pozvati backend za GET /cart
    // kako biste dobili pravi broj stavki pri prvom učitavanju aplikacije
  }

  updateCartCount(newCount: number): void {
    this.cartItemCountSubject.next(newCount);
  }
  clearCartCount(): void {
    this.cartItemCountSubject.next(0);
  }
}
