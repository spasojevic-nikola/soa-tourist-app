import { Component } from '@angular/core';
import { BlogService } from '../blog.service';
import { CreateBlogPayload } from '../model/blog.model';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { Router } from '@angular/router';

@Component({
  selector: 'xp-blog-creation',
  templateUrl: './blog-creation.component.html',
  styleUrls: ['./blog-creation.component.css']
})
export class BlogCreationComponent {

  blogForm: FormGroup;
  imagePreviews: string[] = [];
  isLoading = false;
  errorMessage = '';

  constructor(
    private fb: FormBuilder,
    private blogService: BlogService,
    private router: Router
  ) {
    this.blogForm = this.fb.group({
      title: ['', [Validators.required, Validators.minLength(5)]],
      content: ['', [Validators.required, Validators.minLength(20)]],
      images: [[] as string[]], // Ovde cuvamo Base64 slike
      createdAt: ['', Validators.required],
    });
  }
  onFileSelected(event: Event): void {
    const input = event.target as HTMLInputElement;
    if (!input.files || input.files.length === 0) {
      return;
    }

    const files = Array.from(input.files);
    const imagePromises = files.map(file => this.readFileAsBase64(file));

    Promise.all(imagePromises).then(newBase64images => {
      const currentImages = this.blogForm.get('images')?.value || [];
      const updatedImages = currentImages.concat(newBase64images);
      this.blogForm.patchValue({ images: updatedImages });
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

    const currentImages = this.blogForm.get('images')?.value || [];
    currentImages.splice(indexToRemove, 1);
    this.blogForm.patchValue({ images: currentImages });
  }

  //kreiranje bloga
  onSubmit(): void {
    if (this.blogForm.invalid) {
      this.blogForm.markAllAsTouched();
      return;
    }
    this.isLoading = true;
    this.errorMessage = '';

   
    const payload: CreateBlogPayload = {
      ...this.blogForm.value,
      createdAt: new Date(this.blogForm.value.createdAt)
    };
    

    this.blogService.createBlog(payload).subscribe({
      next: (response) => {
        this.isLoading = false;
        alert('Blog je uspesno kreiran!');
        this.blogForm.reset();
        this.imagePreviews = [];

      },
      error: (err) => {
        this.isLoading = false;
        this.errorMessage = 'Kreiranje bloga nije uspelo. Molimo pokusajte ponovo.';
        console.error(err);
      }
    });
  }
}
