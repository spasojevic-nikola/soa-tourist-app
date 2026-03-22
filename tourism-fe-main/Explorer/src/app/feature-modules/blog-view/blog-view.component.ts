import { Component, OnInit, OnDestroy } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { Subscription } from 'rxjs';
import { BlogService } from '../blog/blog.service';
import { AuthService } from 'src/app/infrastructure/auth/auth.service';
import { StakeholdersService } from 'src/app/infrastructure/stakeholders.service';
import { Blog, BlogComment, AddCommentPayload, UpdateBlogPayload, UpdateCommentPayload } from '../blog/model/blog.model';

@Component({
  selector: 'app-blog-view',
  templateUrl: './blog-view.component.html',
  styleUrls: ['./blog-view.component.css']
})
export class BlogViewComponent implements OnInit, OnDestroy {

  isLoading = true;
  isDetailView = false;
  blogs: Blog[] = [];
  myBlogs: Blog[] = [];
  otherBlogs: Blog[] = [];
  blogDetail?: Blog;
  commentForm!: FormGroup;
  isCommentSending = false;
  currentUsername: string | null = null;

  isEditingBlog = false;
  blogEditForm!: FormGroup;

  editingComment: BlogComment | null = null;
  commentEditForm!: FormGroup;
  isCommentUpdating = false;
  isDeletingBlog = false;


  currentUserId: number | null = null;
  private userSub?: Subscription;
  readonly defaultAvatar = 'assets/images/default_profile.png';
  private authorAvatarCache = new Map<number, string>();
  private pendingAvatarRequests = new Set<number>();

  constructor(
    private blogService: BlogService,
    private route: ActivatedRoute,
    private fb: FormBuilder,
    private authService: AuthService,
    private stakeholdersService: StakeholdersService
  ) {}

  ngOnInit(): void {
    // inicijalizacija forme
    this.commentForm = this.fb.group({
      text: ['', Validators.required]
    });

    this.commentEditForm = this.fb.group({
      editText: ['', Validators.required]
    });

    this.blogEditForm = this.fb.group({
            title: ['', Validators.required],
            content: ['', Validators.required],
        });
      
    // Subscribe na trenutno ulogovanog korisnika
    this.userSub = this.authService.user$.subscribe(user => {
      this.currentUserId = user.id || null;
      this.currentUsername = user.username || null;
      console.log('Current user ID updated:', this.currentUserId, this.currentUsername);
      this.partitionBlogs();
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
    this.blogService.getAllPublishedBlogs().subscribe({
      next: (blogs) => {
        this.blogs = Array.isArray(blogs) ? blogs : [];
        this.blogs.forEach(blog => this.hydrateAuthorAvatar(blog));
        this.partitionBlogs();
        this.isLoading = false;
      },
      error: (err) => console.error(err)
    });
  }

  loadBlogDetail(id: string) {
    this.blogService.getBlogById(id).subscribe({
      next: (blog) => {
        this.blogDetail = blog;
        this.hydrateAuthorAvatar(this.blogDetail);
        this.hydrateCommentList(this.blogDetail.comments);
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
        this.hydrateCommentAvatar(comment);
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
        this.hydrateAuthorAvatar(this.blogDetail);
        this.hydrateCommentList(this.blogDetail.comments);
        this.isLoading = false;
      },
      error: () => this.isLoading = false
    });
  }

   startEditBlog(): void {
        if (!this.blogDetail) return;
        this.isEditingBlog = true;
        this.blogEditForm.patchValue({
            title: this.blogDetail.title,
            content: this.blogDetail.content, 
        });
    }

    cancelEditBlog(): void {
        this.isEditingBlog = false;
        this.blogEditForm.reset();
    }

    saveBlogEdit(): void {
        if (!this.blogDetail || this.blogEditForm.invalid) return;

        const payload: UpdateBlogPayload = {
            title: this.blogEditForm.value.title,
            content: this.blogEditForm.value.content,
            images: this.blogDetail.images, 
        };

        this.blogService.updateBlog(this.blogDetail.id, payload).subscribe({
            next: (updatedBlog) => {
                this.blogDetail = updatedBlog; 
            this.hydrateAuthorAvatar(this.blogDetail);
          this.hydrateCommentList(this.blogDetail?.comments);
                this.isEditingBlog = false;
            },
            error: (err) => console.error('Failed to update blog:', err)
        });
    }
    
    startEditComment(comment: BlogComment): void {
        this.editingComment = comment;
        this.commentEditForm.patchValue({ editText: comment.text });
    }

    cancelEditComment(): void {
        this.editingComment = null;
        this.commentEditForm.reset();
    }

    saveCommentEdit(): void {
        if (!this.blogDetail || !this.editingComment || this.commentEditForm.invalid) return;
        this.isCommentUpdating = true;

        const payload: UpdateCommentPayload = { text: this.commentEditForm.value.editText };
        
        this.blogService.updateComment(this.blogDetail.id, this.editingComment.id, payload).subscribe({
            next: (updatedComment) => {
                const index = this.blogDetail!.comments.findIndex(c => c.id === updatedComment.id);
                if (index !== -1) {
                    this.blogDetail!.comments[index] = updatedComment;
              this.hydrateCommentAvatar(this.blogDetail!.comments[index]);
                }
                
                this.cancelEditComment();
                this.isCommentUpdating = false;
            },
            error: (err) => {
                console.error('Failed to update comment:', err);
                this.isCommentUpdating = false;
            }
        });
    }

      deleteBlog(): void {
        if (!this.blogDetail) {
          return;
        }

        const confirmed = window.confirm('Are you sure you want to delete this blog? This action cannot be undone.');
        if (!confirmed) {
          return;
        }

        this.isDeletingBlog = true;
        const blogId = this.blogDetail.id;

        this.blogService.deleteBlog(blogId).subscribe({
          next: () => {
            this.blogs = this.blogs.filter(blog => blog.id !== blogId);
            this.partitionBlogs();
            this.blogDetail = undefined;
            this.isDetailView = false;
            this.isDeletingBlog = false;
          },
          error: (err) => {
            console.error('Failed to delete blog:', err);
            this.isDeletingBlog = false;
          }
        });
      }

    private hydrateAuthorAvatar(blog?: Blog): void {
      if (!blog) {
        return;
      }

      this.assignAvatar(blog.authorId, avatar => {
        blog.authorProfileImage = avatar;
      });
    }

    private hydrateCommentList(comments?: BlogComment[]): void {
      if (!comments || comments.length === 0) {
        return;
      }

      comments.forEach(comment => this.hydrateCommentAvatar(comment));
    }

    private hydrateCommentAvatar(comment?: BlogComment): void {
      if (!comment) {
        return;
      }

      this.assignAvatar(comment.authorId, avatar => {
        comment.authorProfileImage = avatar;
      });
    }

    private assignAvatar(authorId: number | null | undefined, apply: (avatar: string) => void): void {
      if (authorId === null || authorId === undefined) {
        return;
      }

      const cachedAvatar = this.authorAvatarCache.get(authorId);
      if (cachedAvatar) {
        apply(cachedAvatar);
        return;
      }

      apply(this.defaultAvatar);

      if (this.pendingAvatarRequests.has(authorId)) {
        return;
      }

      this.pendingAvatarRequests.add(authorId);

      this.stakeholdersService.getUserById(authorId).subscribe({
        next: (user) => {
          const avatarUrl = user.profile_image || this.defaultAvatar;
          this.authorAvatarCache.set(authorId, avatarUrl);
          this.updateAvatarAssignments(authorId, avatarUrl);
          this.pendingAvatarRequests.delete(authorId);
        },
        error: () => {
          this.authorAvatarCache.set(authorId, this.defaultAvatar);
          this.updateAvatarAssignments(authorId, this.defaultAvatar);
          this.pendingAvatarRequests.delete(authorId);
        }
      });
    }

    private updateAvatarAssignments(authorId: number, avatarUrl: string): void {
      this.blogs.forEach(blog => {
        if (blog.authorId === authorId) {
          blog.authorProfileImage = avatarUrl;
        }

        blog.comments?.forEach(comment => {
          if (comment.authorId === authorId) {
            comment.authorProfileImage = avatarUrl;
          }
        });
      });

      if (this.blogDetail && this.blogDetail.authorId === authorId) {
        this.blogDetail.authorProfileImage = avatarUrl;
      }

      if (this.blogDetail?.comments) {
        this.blogDetail.comments.forEach(comment => {
          if (comment.authorId === authorId) {
            comment.authorProfileImage = avatarUrl;
          }
        });
      }
    }

    private partitionBlogs(): void {
      if (!this.blogs || this.blogs.length === 0) {
        this.myBlogs = [];
        this.otherBlogs = [];
        return;
      }

      if (this.currentUserId === null || this.currentUserId === undefined) {
        this.myBlogs = [];
        this.otherBlogs = [...this.blogs];
        return;
      }

      this.myBlogs = this.blogs.filter(blog => blog.authorId === this.currentUserId);
      this.otherBlogs = this.blogs.filter(blog => blog.authorId !== this.currentUserId);
    }
}
