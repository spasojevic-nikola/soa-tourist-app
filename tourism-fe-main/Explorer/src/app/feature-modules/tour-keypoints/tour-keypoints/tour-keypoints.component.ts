import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { KeypointService } from '../keypoint.service';
import { CreateKeyPointPayload, KeyPoint } from '../model/keypoint.model';

@Component({
  selector: 'xp-tour-keypoints',
  templateUrl: './tour-keypoints.component.html',
  styleUrls: ['./tour-keypoints.component.css']
})
export class TourKeypointsComponent implements OnInit {
  tourId: number;
  keyPointsForm: FormGroup;
  isSubmitting = false;
  successMessage: string | null = null;
  errorMessage: string | null = null;
  imagePreview: string | null = null;  
  addedKeyPoints: KeyPoint[] = [];

  constructor(
    private route: ActivatedRoute,
    private fb: FormBuilder,
    private keyPointService: KeypointService
  ) {}

  ngOnInit(): void {
    this.tourId = +this.route.snapshot.paramMap.get('tourId')!;
    this.initForm();
  }

  initForm(): void {
    this.keyPointsForm = this.fb.group({
      name: ['', [Validators.required, Validators.minLength(3)]],
      description: ['', [Validators.required, Validators.minLength(10)]],
      latitude: ['', [Validators.required, Validators.min(-90), Validators.max(90)]],
      longitude: ['', [Validators.required, Validators.min(-180), Validators.max(180)]],
      image: [''], 
      order: ['', [Validators.required, Validators.min(1)]]
    });
  }

  // METODA ZA BASE64 CONVERSION
  onFileSelected(event: any): void {
    const file = event.target.files[0];
    if (file) {
      if (!file.type.match('image.*')) {
        this.errorMessage = 'Please select a valid image file.';
        return;
      }

      const reader = new FileReader();
      reader.onload = () => {
        this.imagePreview = reader.result as string;
        this.keyPointsForm.patchValue({
          image: reader.result as string
        });
      };
      reader.readAsDataURL(file);
    }
  }

  removeImage(): void {
    this.imagePreview = null;
    this.keyPointsForm.patchValue({ image: '' });
  }

  onSubmit(): void {
    if (this.keyPointsForm.invalid) {
      this.keyPointsForm.markAllAsTouched();
      return;
    }

    this.isSubmitting = true;
    this.successMessage = null;
    this.errorMessage = null;

    const formValue = this.keyPointsForm.value;
    const payload: CreateKeyPointPayload = {
      name: formValue.name,
      description: formValue.description,
      latitude: parseFloat(formValue.latitude),
      longitude: parseFloat(formValue.longitude),
      image: formValue.image, // BASE64 STRING
      order: parseInt(formValue.order)
    };

    this.keyPointService.createKeyPoint(this.tourId, payload).subscribe({
      next: (createdKeyPoint) => {
        this.addedKeyPoints.push(createdKeyPoint);
        this.successMessage = `Key point "${createdKeyPoint.name}" added successfully!`;
        this.keyPointsForm.reset();
        this.imagePreview = null;
        this.isSubmitting = false;
      },
      error: (err) => {
        this.errorMessage = `Error: ${err.error?.error || err.error || 'Please check your data and try again.'}`;
        console.error(err);
        this.isSubmitting = false;
      }
    });
  }

  addAnother(): void {
    this.successMessage = null;
    this.keyPointsForm.reset();
    this.imagePreview = null;
  }

  finish(): void {
    window.history.back();
  }

  removeKeyPoint(keyPointId: number, index: number): void {
    this.keyPointService.deleteKeyPoint(keyPointId).subscribe({
      next: () => {
        this.addedKeyPoints.splice(index, 1);
        this.successMessage = 'Key point removed successfully!';
      },
      error: (err) => {
        this.errorMessage = `Error removing key point: ${err.error?.error || err.error || 'Please try again.'}`;
      }
    });
  }
}