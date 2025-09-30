import { Component, OnInit, Inject } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { CreateKeyPointPayload } from '../model/keypoint.model';

export interface KeypointDialogData {
  latitude: number;
  longitude: number;
  order: number;
  existingKeyPoint?: CreateKeyPointPayload;
}

@Component({
  selector: 'xp-tour-keypoints',
  templateUrl: './tour-keypoints.component.html',
  styleUrls: ['./tour-keypoints.component.css']
})
export class TourKeypointsComponent implements OnInit {
  keyPointsForm: FormGroup;
  isSubmitting = false;
  successMessage: string | null = null;
  errorMessage: string | null = null;
  imagePreview: string | null = null;
  isEditMode = false;

  constructor(
    private fb: FormBuilder,
    public dialogRef: MatDialogRef<TourKeypointsComponent>,
    @Inject(MAT_DIALOG_DATA) public data: KeypointDialogData
  ) {}

  ngOnInit(): void {
    this.initForm();
    this.patchFormWithData();

    // Ako postoje postojeći podaci, popuni formu
  if (this.data.existingKeyPoint) {
    this.patchFormWithExistingData();
  }
}

patchFormWithExistingData(): void {
  const existing = this.data.existingKeyPoint;
  if(existing != null)
  {
    this.keyPointsForm.patchValue({
    name: existing.name,
    description: existing.description,
    latitude: existing.latitude,
    longitude: existing.longitude,
    order: existing.order
  });

    // Ako postoji slika, postavi preview
    if (existing.image && typeof existing.image === 'string') {
      this.imagePreview = existing.image;
    }

    this.isEditMode = true;
  }
  
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

  patchFormWithData(): void {
    this.keyPointsForm.patchValue({
      latitude: this.data.latitude,
      longitude: this.data.longitude,
      order: this.data.order
    });
  }

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
      image: formValue.image,
      order: parseInt(formValue.order)
    };

    // Zatvori dialog i vrati podatke
    this.dialogRef.close(payload);
  }

  onCancel(): void {
    this.dialogRef.close(); // Zatvori bez podataka
  }

}