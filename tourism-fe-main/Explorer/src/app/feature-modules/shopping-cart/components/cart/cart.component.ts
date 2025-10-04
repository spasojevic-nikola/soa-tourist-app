import { Component, OnInit } from '@angular/core';
import { ShoppingCart } from '../../services/models/cart.model';
import { CartService } from '../../services/cart.service';
import { PurchaseTokenResponse } from '../../services/dto/purchase-token-response.dto';
import { MatSnackBar } from '@angular/material/snack-bar';
import { CartStateService } from '../../services/cart-state.service';
import { Router } from '@angular/router';

@Component({
  selector: 'xp-cart',
  templateUrl: './cart.component.html',
  styleUrls: ['./cart.component.css']
})
export class CartComponent implements OnInit {

  cart: ShoppingCart | null = null;
  isLoading = true;
  checkoutInProgress = false;

  constructor(
    private cartService: CartService,
    private snackBar: MatSnackBar,
    private cartStateService : CartStateService,
    private router: Router
  ) { }

  ngOnInit(): void {
    // 💡 Pri učitavanju komponente, dohvati trenutno stanje korpe
    this.loadCart();
  }

  loadCart(): void {
    this.isLoading = true;
    this.cartService.getCart().subscribe({
      next: (data) => {
        // Uspesno ucitavanje podataka
        this.cart = data;
        this.isLoading = false;
      },
      error: (err) => {
        console.error('Failed to load cart:', err);
        this.snackBar.open('Error loading cart contents. Log in required.', 'Dismiss', { duration: 4000 });
        this.isLoading = false;
      }
    });
  }

  onCheckout(): void {
    if (!this.cart || this.cart.items.length === 0) {
      this.snackBar.open('Cart is empty. Add tours to proceed.', 'Dismiss', { duration: 3000 });
      return;
    }

    this.checkoutInProgress = true;
    
    // Poziv Backendu
    this.cartService.checkout().subscribe({
      next: (response: PurchaseTokenResponse) => {
        // Obavesti korisnika o uspehu i prikaži poruku iz Backenda
        this.snackBar.open(response.message, 'Success!', { duration: 6000 });
        this.checkoutInProgress = false;
        
        // Ažuriraj lokalno stanje na praznu korpu (UI odmah reaguje)
        this.cart = { ...this.cart!, items: [], total: 0 }; 
      },
      error: (err) => {
        console.error('Checkout failed:', err);
        this.snackBar.open('Checkout failed. Please try again.', 'Dismiss', { duration: 5000 });
        this.checkoutInProgress = false;
      }
    });
  }

  onRemoveItem(tourId: string, itemName: string): void {
    this.cartService.removeItem(tourId).subscribe({
        next: (updatedCart) => {
            // Ažuriraj lokalno stanje na novu korpu
            this.cart = updatedCart; 
            this.snackBar.open(`${itemName} successfully removed.`, 'Dismiss', { duration: 3000 });

            this.cartStateService.updateCartCount(updatedCart.items.length);
        },
        error: (err) => {
            console.error('Failed to remove item:', err);
            this.snackBar.open('Failed to remove item. Check console.', 'Dismiss', { duration: 4000 });
        }
    });
}
navigateToTour(tourId: string): void {
  this.router.navigate(['/tours', tourId]);
}
}
