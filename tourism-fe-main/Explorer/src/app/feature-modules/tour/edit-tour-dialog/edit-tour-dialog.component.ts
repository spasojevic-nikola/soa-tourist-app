import { Component, Inject, OnInit } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { Tour } from '../model/tour.model';

@Component({
  selector: 'xp-edit-tour-dialog',
  templateUrl: './edit-tour-dialog.component.html',
  styleUrls: ['./edit-tour-dialog.component.css']
})
export class EditTourDialogComponent implements OnInit {
  editForm: FormGroup;

  constructor(
    private fb: FormBuilder,
    public dialogRef: MatDialogRef<EditTourDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public data: { tour: Tour }
  ) {
    this.editForm = this.fb.group({
      name: [data.tour.name, [Validators.required, Validators.minLength(5)]],
      description: [data.tour.description, [Validators.required, Validators.minLength(20)]],
      difficulty: [data.tour.difficulty, Validators.required],
      price: [data.tour.price || 0, [Validators.required, Validators.min(0)]],
      tags: [data.tour.tags ? data.tour.tags.join(', ') : '', Validators.required]
    });
  }

  ngOnInit(): void {}

  onCancel(): void {
    this.dialogRef.close();
  }

  onSubmit(): void {
    if (this.editForm.valid) {
      const formValue = this.editForm.value;
      const tagsArray = formValue.tags.split(',')
        .map((t: string) => t.trim())
        .filter((t: string) => t !== '');

      this.dialogRef.close({
        name: formValue.name,
        description: formValue.description,
        difficulty: formValue.difficulty,
        price: formValue.price,
        tags: tagsArray
      });
    }
  }
}
