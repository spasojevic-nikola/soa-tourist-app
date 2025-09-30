import { Component } from '@angular/core';
import { User } from '../profile/model/profile.model';
import { ActivatedRoute } from '@angular/router';
import { StakeholdersService } from 'src/app/infrastructure/stakeholders.service';
import { FollowerService } from 'src/app/infrastructure/follower/follower.service';
import { User as AuthUser } from 'src/app/infrastructure/auth/model/user.model';
import { AuthService } from 'src/app/infrastructure/auth/auth.service';

@Component({
  selector: 'xp-view-profile',
  templateUrl: './view-profile.component.html',
  styleUrls: ['./view-profile.component.css']
})
export class ViewProfileComponent {
  user: User;
  isLoading = true;
  isFollowing: boolean = false;
  isMyProfile: boolean = false;
  loggedInUser: AuthUser;


  constructor(
    private route: ActivatedRoute,
    private stakeholdersService: StakeholdersService,
    private followerService: FollowerService,
    private authService: AuthService
  ) { }

  ngOnInit(): void {
    // Prvo dobavljamo podatke o ulogovanom korisniku
    this.authService.user$.subscribe(user => {
      this.loggedInUser = user;
    });

    // Zatim slušamo promene u ruti
    this.route.params.subscribe(params => {
      const userIdToView = params['id'];

      // Proveravamo da li korisnik gleda svoj profil
      if (this.loggedInUser && this.loggedInUser.id == userIdToView) {
        this.isMyProfile = true;
      }
      
      this.loadUserProfile(userIdToView);
    });
  }

  // Funkcija za učitavanje profila
  loadUserProfile(userId: number): void {
    this.isLoading = true;
    this.stakeholdersService.getUserById(userId).subscribe(user => {
      this.user = user;
      
      // Ako ne gledamo svoj profil, proveravamo da li ga pratimo
      if (!this.isMyProfile) {
        this.checkFollowingStatus();
      }
      
      this.isLoading = false;
    });
  }

  // Proverava status praćenja
  checkFollowingStatus(): void {
    this.followerService.checkIfFollowing(this.user.id).subscribe(result => {
      this.isFollowing = result.follows;
    });
  }


  follow(): void {
    this.followerService.follow(this.user.id).subscribe({
      next: () => {
        this.isFollowing = true;
      }
    });
  }

  unfollow(): void {
    this.followerService.unfollow(this.user.id).subscribe({
      next: () => {
        this.isFollowing = false;
      }
    });
  }
}