import { Component, Inject, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { Review } from '../model/review.model';

export interface ReviewDialogData {
  tourId: number;
  tourName: string;
  review?: Review; // If editing existing review
}

@Component({
  selector: 'xp-review-dialog',
  templateUrl: './review-dialog.component.html',
  styleUrls: ['./review-dialog.component.css']
})
export class ReviewDialogComponent implements OnInit {
  reviewForm: FormGroup;
  isEditMode = false;
  maxDate = new Date();
  imagePreviews: string[] = [];

  constructor(
    private fb: FormBuilder,
    public dialogRef: MatDialogRef<ReviewDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public data: ReviewDialogData
  ) {
    this.isEditMode = !!data.review;
    
    this.reviewForm = this.fb.group({
      rating: [data.review?.rating || 5, [Validators.required, Validators.min(1), Validators.max(5)]],
      comment: [data.review?.comment || '', Validators.required],
      visitDate: [data.review?.visitDate ? new Date(data.review.visitDate) : new Date(), Validators.required],
      images: [data.review?.images || []]
    });

    // Load existing images for preview if editing
    if (data.review?.images) {
      this.imagePreviews = [...data.review.images];
    }
  }

  ngOnInit(): void {}

  onFileSelected(event: Event): void {
    const input = event.target as HTMLInputElement;
    if (!input.files || input.files.length === 0) {
      return;
    }

    const files = Array.from(input.files);
    const imagePromises = files.map(file => this.readFileAsBase64(file));

    Promise.all(imagePromises).then(newBase64images => {
      const currentImages = this.reviewForm.get('images')?.value || [];
      const updatedImages = currentImages.concat(newBase64images);
      this.reviewForm.patchValue({ images: updatedImages });
    });

    input.value = '';
  }

  private readFileAsBase64(file: File): Promise<string> {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();
      reader.onload = () => {
        this.imagePreviews.push(reader.result as string);
        resolve(reader.result as string);
      };
      reader.onerror = (error) => reject(error);
      reader.readAsDataURL(file);
    });
  }

  removeImage(indexToRemove: number): void {
    this.imagePreviews.splice(indexToRemove, 1);

    const currentImages = this.reviewForm.get('images')?.value || [];
    currentImages.splice(indexToRemove, 1);
    this.reviewForm.patchValue({ images: currentImages });
  }

  onCancel(): void {
    this.dialogRef.close();
  }

  onSubmit(): void {
    if (this.reviewForm.valid) {
      const formValue = this.reviewForm.value;
      
      // Convert visitDate to ISO string if it's a Date object
      let visitDateISO: string;
      if (formValue.visitDate instanceof Date) {
        visitDateISO = formValue.visitDate.toISOString();
      } else if (typeof formValue.visitDate === 'string') {
        visitDateISO = new Date(formValue.visitDate).toISOString();
      } else {
        visitDateISO = new Date().toISOString();
      }
      
      const result = {
        tourId: this.data.tourId,
        rating: formValue.rating,
        comment: formValue.comment,
        visitDate: visitDateISO,
        images: formValue.images
      };
      this.dialogRef.close(result);
    }
  }

  get ratingArray(): number[] {
    return [1, 2, 3, 4, 5];
  }
}
