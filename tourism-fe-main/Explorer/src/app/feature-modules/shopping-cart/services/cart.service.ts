import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { environment } from 'src/env/environment';
import { ShoppingCart } from './models/cart.model';
import { AddItemPayload } from './dto/add-item.dto';
import { PurchaseTokenResponse } from './dto/purchase-token-response.dto';

@Injectable({
  providedIn: 'root'
})
export class CartService {

  private apiUrl = `${environment.purchaseApiHost}/cart`; 

  constructor(private http: HttpClient) { }
  getCart(): Observable<ShoppingCart> {
    return this.http.get<ShoppingCart>(this.apiUrl);
  }
  
  addItem(payload: AddItemPayload): Observable<ShoppingCart> {
    //  salje AddItemPayload
    return this.http.post<ShoppingCart>(`${this.apiUrl}/items`, payload); 
  }

  checkout(): Observable<PurchaseTokenResponse> {
    //  ocekuje PurchaseTokenResponse
    return this.http.post<PurchaseTokenResponse>(`${this.apiUrl}/checkout`, {});
  }
  
  removeItem(tourId: string): Observable<ShoppingCart> {
    const url = `${this.apiUrl}/items/${tourId}`;
    return this.http.delete<ShoppingCart>(url);
  }
}
