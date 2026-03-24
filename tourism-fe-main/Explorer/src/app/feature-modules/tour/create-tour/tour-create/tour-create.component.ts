import { Component, OnInit, Output, EventEmitter } from '@angular/core'; // DODAJ Output, EventEmitter
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { TourService } from '../../tour.service';
import { CreateTourPayload } from '../../dto/tour-creation.dto';
import { Router } from '@angular/router';

@Component({
  selector: 'xp-tour-create',
  templateUrl: './tour-create.component.html',
  styleUrls: ['./tour-create.component.css']
})
export class TourCreateComponent implements OnInit {

  @Output() tourCreated = new EventEmitter<CreateTourPayload>(); // DODAJ OVO

  tourForm: FormGroup;
  isSubmitting = false;
  successMessage: string | null = null;
  errorMessage: string | null = null;

  showKeyPointDialog = false;  
  createdTourId: number | null = null;  

  constructor(
    private fb: FormBuilder,
    private tourService: TourService,
    private router: Router
  ) { }

  ngOnInit(): void {
    this.tourForm = this.fb.group({
      name: ['', [Validators.required, Validators.minLength(5)]],
      description: ['', [Validators.required, Validators.minLength(20)]],
      difficulty: ['Easy', Validators.required], 
      price: [0, [Validators.required, Validators.min(0)]],
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

    const tagsArray = formValue.tags.split(',')
      .map((tag: string) => tag.trim())
      .filter((tag: string) => tag !== '');

    const payload: CreateTourPayload = {
      name: formValue.name,
      description: formValue.description,
      difficulty: formValue.difficulty,
      price: formValue.price,
      tags: tagsArray
    };
    
    this.tourCreated.emit(payload); // EMITUJ PAYLOAD WIZARDU
    this.isSubmitting = false;

  }

}
