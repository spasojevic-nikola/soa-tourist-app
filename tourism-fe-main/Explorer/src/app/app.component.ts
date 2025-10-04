import { Component, OnInit } from '@angular/core';
import { AuthService } from './infrastructure/auth/auth.service';
import { CartService } from './feature-modules/shopping-cart/services/cart.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent implements OnInit {
  title = 'Explorer';

  constructor(
    private authService: AuthService,
    private cartService: CartService
  ) {}


  ngOnInit(): void {
    this.checkIfUserExists();
    this.authService.user$.subscribe(user => {
      if (user && user.username !== '') {
        this.cartService.getCart().subscribe(); // Automatski ažurira count kroz tap()
      }
    });
    
  }
  
  private checkIfUserExists(): void {
    this.authService.checkIfUserExists();
  }
}
