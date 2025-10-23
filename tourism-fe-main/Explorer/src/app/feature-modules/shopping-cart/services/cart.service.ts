import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable, tap } from 'rxjs';
import { environment } from 'src/env/environment';
import { ShoppingCart } from './models/cart.model';
import { AddItemPayload } from './dto/add-item.dto';
import { PurchaseTokenResponse } from './dto/purchase-token-response.dto';
import { CartStateService } from './cart-state.service';

@Injectable({
  providedIn: 'root'
})
export class CartService {

  private apiUrl = `${environment.purchaseApiHost}`; 

  constructor(private http: HttpClient,
    private cartStateService: CartStateService

  ) { }
  getCart(): Observable<ShoppingCart> {
    return this.http.get<ShoppingCart>(this.apiUrl).pipe(
      tap(cart => this.cartStateService.updateCartCount(cart.items.length))
    );
    }
  
  addItem(payload: AddItemPayload): Observable<ShoppingCart> {
    //  salje AddItemPayload
    return this.http.post<ShoppingCart>(`${this.apiUrl}/items`, payload).pipe(
      tap(cart => this.cartStateService.updateCartCount(cart.items.length))
    );  }

  checkout(): Observable<PurchaseTokenResponse> {
    //  ocekuje PurchaseTokenResponse
    return this.http.post<PurchaseTokenResponse>(`${this.apiUrl}/checkout`, {}).pipe(
      tap(() => this.cartStateService.updateCartCount(0))
    );  }
  
  removeItem(tourId: string): Observable<ShoppingCart> {
    const url = `${this.apiUrl}/items/${tourId}`;
    return this.http.delete<ShoppingCart>(url).pipe(
      tap(cart => this.cartStateService.updateCartCount(cart.items.length))
    );
  }
  hasPurchased(tourId: string): Observable<{ isPurchased: boolean }> {
    return this.http.get<{ isPurchased: boolean }>(`${this.apiUrl}/purchase-status/${tourId}`);
  }
}
