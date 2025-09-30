import { Component, OnInit } from '@angular/core';
import { FormControl } from '@angular/forms';
import { Router } from '@angular/router';
import { debounceTime, distinctUntilChanged, filter, Observable, of, switchMap } from 'rxjs';
import { User as AuthUser } from 'src/app/infrastructure/auth/model/user.model';
import { AuthService } from 'src/app/infrastructure/auth/auth.service';
import { User } from 'src/app/feature-modules/user-profile/profile/model/profile.model';
import { StakeholdersService } from 'src/app/infrastructure/stakeholders.service';

@Component({
  selector: 'xp-navbar',
  templateUrl: './navbar.component.html',
  styleUrls: ['./navbar.component.css']
})
export class NavbarComponent implements OnInit {

  user: AuthUser | undefined;
  searchControl = new FormControl('');
  filteredUsers$: Observable<User[]>;

  constructor(private authService: AuthService,    
    private stakeholdersService: StakeholdersService,
    private router: Router) {}

    ngOnInit(): void {
      this.authService.user$.subscribe(user => {
        this.user = user;
      });
  
      // Logika za "live search"
      this.filteredUsers$ = this.searchControl.valueChanges.pipe(
        // Sacekaj 300ms nakon sto korisnik prestane da kuca
        debounceTime(300),
        // Ne salji zahtev ako je tekst isti kao prethodni
        distinctUntilChanged(),
        // Ne pretrazuj ako je unos kraci od 2 karaktera
        filter(query => typeof query === 'string' && query.length > 1),
        // Posalji zahtev i otkaži sve prethodne ako stigne novi unos
        switchMap(query => {
          if (query) {
            return this.stakeholdersService.searchUsers(query);
          } else {
            return of([]); // Vrati prazan niz ako je polje prazno
          }
        })
      );
    }

  onLogout(): void {
    this.authService.logout();
  }
  onSearch(query: string): void {
    if (query.trim()) { // Proveravamo da string nije prazan
      // Preusmeravamo korisnika na /search stranicu sa query parametrom 'q'
      this.router.navigate(['/search'], { queryParams: { q: query } });
    }
  }
}
