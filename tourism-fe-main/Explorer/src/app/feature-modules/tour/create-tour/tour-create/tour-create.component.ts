import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { TourService } from '../../tour.service';
import { CreateTourPayload } from '../../dto/tour-creation.dto';

@Component({
  selector: 'xp-tour-create',
  templateUrl: './tour-create.component.html',
  styleUrls: ['./tour-create.component.css']
})
export class TourCreateComponent implements OnInit {

  tourForm: FormGroup;
  isSubmitting = false;
  successMessage: string | null = null;
  errorMessage: string | null = null;

  constructor(
    private fb: FormBuilder,
    private tourService: TourService
  ) { }

  ngOnInit(): void {
    this.tourForm = this.fb.group({
      name: ['', [Validators.required, Validators.minLength(5)]],
      description: ['', [Validators.required, Validators.minLength(20)]],
      difficulty: ['Easy', Validators.required], 
      tags: ['', Validators.required]         
    });
  }

  onSubmit(): void {
    if (this.tourForm.invalid) {
      this.tourForm.markAllAsTouched();
      return;
    }

    this.isSubmitting = true;
    this.successMessage = null;
    this.errorMessage = null;

    const formValue = this.tourForm.value;
    
    // Tagove koje korisnik unese kao string ("planina, reka") pretvaramo u niz ["planina", "reka"]
    const tagsArray = formValue.tags.split(',')
      .map((tag: string) => tag.trim())
      .filter((tag: string) => tag !== ''); 

    // Kreiramo DTO koji saljemo
    const payload: CreateTourPayload = {
      name: formValue.name,
      description: formValue.description,
      difficulty: formValue.difficulty,
      tags: tagsArray
    };

    this.tourService.createTour(payload).subscribe({
      next: (createdTour) => {
       // this.successMessage = `Tura "${createdTour.name}" je uspešno kreirana i ima status 'Draft'.`;
        this.tourForm.reset({ difficulty: 'Easy', name: '', description: '', tags: '' }); // Resetujemo formu
        this.isSubmitting = false;
      },
      error: (err) => {
        this.errorMessage = `Došlo je do greške: ${err.error.error || err.error || 'Proverite podatke i pokušajte ponovo.'}`;
        console.error(err);
        this.isSubmitting = false;
      }
    });
  }
}