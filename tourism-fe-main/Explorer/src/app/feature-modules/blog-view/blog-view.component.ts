import { Component, OnInit, OnDestroy } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { Subscription } from 'rxjs';
import { BlogService } from '../blog/blog.service';
import { AuthService } from 'src/app/infrastructure/auth/auth.service';
import { TokenStorage } from 'src/app/infrastructure/auth/jwt/token.service';
import { Blog, BlogComment, AddCommentPayload } from '../blog/model/blog.model';

@Component({
  selector: 'app-blog-view',
  templateUrl: './blog-view.component.html',
  styleUrls: ['./blog-view.component.css']
})
export class BlogViewComponent implements OnInit, OnDestroy {

  isLoading = true;
  isDetailView = false;
  blogs: Blog[] = [];
  blogDetail?: Blog;
  commentForm!: FormGroup;
  isCommentSending = false;
  currentUsername: string | null = null;


  currentUserId: number | null = null;
  private userSub?: Subscription;

  constructor(
    private blogService: BlogService,
    private route: ActivatedRoute,
    private fb: FormBuilder,
    private authService: AuthService
  ) {}

  ngOnInit(): void {
    // inicijalizacija forme
    this.commentForm = this.fb.group({
      text: ['', Validators.required]
    });

    // Subscribe na trenutno ulogovanog korisnika
    this.userSub = this.authService.user$.subscribe(user => {
      this.currentUserId = user.id || null;
      this.currentUsername = user.username || null;
      console.log('Current user ID updated:', this.currentUserId, this.currentUsername);
    });

    const blogId = this.route.snapshot.paramMap.get('id');
    if (blogId) {
      this.isDetailView = true;
      this.loadBlogDetail(blogId);
    } else {
      this.loadAllBlogs();
    }
  }

  ngOnDestroy(): void {
    this.userSub?.unsubscribe();
  }

  loadAllBlogs() {
    this.blogService.getAllBlogs().subscribe({
      next: (blogs) => {
        this.blogs = blogs;
        this.isLoading = false;
      },
      error: (err) => console.error(err)
    });
  }

  loadBlogDetail(id: string) {
    this.blogService.getBlogById(id).subscribe({
      next: (blog) => {
        this.blogDetail = blog;
        this.isLoading = false;
      },
      error: (err) => console.error(err)
    });
  }

  addComment() {
    if (!this.blogDetail) return;
    this.isCommentSending = true;

    const payload: AddCommentPayload = { text: this.commentForm.value.text };
    this.blogService.addComment(this.blogDetail.id, payload).subscribe({
      next: (comment: BlogComment) => {
      comment.authorUsername = this.currentUsername ?? 'Unknown';
        this.blogDetail?.comments.push(comment);
        this.commentForm.reset();
        this.isCommentSending = false;
      },
      error: (err) => {
        console.error(err);
        this.isCommentSending = false;
      }
    });
  }

  toggleLike() {
    if (!this.blogDetail || !this.currentUserId) return;

    this.blogService.toggleLike(this.blogDetail.id).subscribe({
      next: () => {
        const index = this.blogDetail!.likes.indexOf(this.currentUserId!);
        if (index === -1) {
          this.blogDetail!.likes.push(this.currentUserId!);
        } else {
          this.blogDetail!.likes.splice(index, 1);
        }
      },
      error: (err) => console.error(err)
    });
  }

  goToDetail(blogId: string) {
    this.isDetailView = true;
    this.isLoading = true;
    this.blogService.getBlogById(blogId).subscribe({
      next: (data) => {
        this.blogDetail = data;
        this.isLoading = false;
      },
      error: () => this.isLoading = false
    });
  }
}
