import { Component, OnInit } from '@angular/core';
import { ProfileService } from '../profile.service';
import { User } from './model/profile.model';

@Component({
  selector: 'app-profile',
  templateUrl: './profile.component.html',
  styleUrls: ['./profile.component.css']
})
export class ProfileComponent implements OnInit {
  user: User | undefined;
  isLoading = true;

  constructor(private profileService: ProfileService) { }

  ngOnInit(): void {
    this.fetchUserProfile();
  }

  fetchUserProfile(): void {
    this.profileService.getProfile().subscribe({
      next: (data) => {
        this.user = data;
        this.isLoading = false;
        console.log('User profile fetched successfully', this.user);
      },
      error: (err) => {
        console.error('Failed to fetch user profile', err);
        this.isLoading = false;
        // Možete dodati logiku za obradu greške, npr. preusmeravanje na login
      }
    });
  }

}