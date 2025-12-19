import { Component, OnInit } from '@angular/core';
import { ProfileService } from '../profile.service';
import { UpdateUserProfilePayload, User } from './model/profile.model';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';

@Component({
  selector: 'app-profile',
  templateUrl: './profile.component.html',
  styleUrls: ['./profile.component.css']
})
export class ProfileComponent implements OnInit {
  user: User | undefined;
  isLoading = true;

  profileForm: FormGroup; // FormGroup koji će kontrolisati  formu
  isEditing = false; // Prekidac za prikazivanje forme za izmenu
  imagePreviewUrl: string | ArrayBuffer | null = null; //za privremeni prikaz slike koju korisnik odabere

  errorMessage = '';

  constructor(private profileService: ProfileService,
              private fb: FormBuilder // Iza lakse kreiranje forme

            ) { 
                //inicijalizujemo formu
                this.profileForm = this.fb.group({
                // Definišemo polja forme i njihove validatore
                firstName: ['', Validators.required],
                lastName: ['', Validators.required],
                profileImageFile: [null], 
                biography: ['', Validators.maxLength(500)],
                motto: ['', Validators.maxLength(100)]
    });
            }

  ngOnInit(): void {
    this.fetchUserProfile();
  }

  fetchUserProfile(): void {
    this.profileService.getProfile().subscribe({
      next: (data) => {
        this.user = data;
        this.isLoading = false;

        this.profileForm.patchValue({
          firstName: data.first_name,
          lastName: data.last_name,
          biography: data.biography,
          motto: data.motto,
        });

        console.log('User profile fetched successfully', this.user);
      },
      error: (err) => {
        console.error('Failed to fetch user profile', err);
        this.isLoading = false;
        // Možete dodati logiku za obradu greške, npr. preusmeravanje na login
      }
    });
  }

  saveProfile(): void {
    if (this.profileForm.invalid) {
      // Ako forma nije validna ne radimo nis
      return;
    }

    this.isLoading = true;
    this.errorMessage = '';

    // Kreiramo 'payload' objekat sa podacima iz forme
    const formValues = this.profileForm.value;
    const imagePayload = this.imagePreviewUrl ? this.imagePreviewUrl.toString() : this.user?.profile_image;

    const payload: UpdateUserProfilePayload = {
      first_name: formValues.firstName,
      last_name: formValues.lastName,
      profile_image: imagePayload, // Šaljemo Base64 string ili stari URL
      biography: formValues.biography,
      motto: formValues.motto,
    };

    // Pozivamo servis da pošalje podatke na backend
    this.profileService.updateProfile(payload).subscribe({
      next: (updatedUser) => {
        this.user = updatedUser; // Ažuriramo prikaz novim podacima
        this.isEditing = false; // Vraćamo se na prikaz profila (gasimo formu)
        this.isLoading = false;
        this.imagePreviewUrl = null; // Resetujemo preview
        //alert('Profile updated successfully!');
      },
      error: (err) => {
        this.errorMessage = 'Failed to update profile. Please try again.';
        console.error(err);
        this.isLoading = false;
      }
    });
  }


  enterEditMode(): void {
    this.isEditing = true;
  }

  cancelEdit(): void {
    this.isEditing = false;
    this.imagePreviewUrl = null; // Resetujemo preview i pri odustajanju
    // Vraćamo vrednosti forme na originalne
    if (this.user) {
      this.profileForm.patchValue({
        firstName: this.user.first_name,
        lastName: this.user.last_name,
        profileImage: this.user.profile_image,
        biography: this.user.biography,
        motto: this.user.motto,
      });
    }
  }

  onFileSelected(event: Event): void {
    const input = event.target as HTMLInputElement;
    if (input.files && input.files[0]) {
      const file = input.files[0];
      
      // Sačuvamo fajl u formu
      this.profileForm.patchValue({ profileImageFile: file });
      
      // Kreiramo preview slike
      const reader = new FileReader();
      reader.onload = () => {
        this.imagePreviewUrl = reader.result;
      };
      reader.readAsDataURL(file);
    }
  }
}