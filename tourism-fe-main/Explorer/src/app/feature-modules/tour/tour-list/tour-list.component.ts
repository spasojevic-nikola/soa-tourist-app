import { Component, OnDestroy, OnInit } from '@angular/core';
import { Observable, Subscription } from 'rxjs';
import { tap } from 'rxjs/operators';
import { TourService } from '../tour.service';
import { Tour } from '../model/tour.model';
import { AuthService } from 'src/app/infrastructure/auth/auth.service';
import { Router } from '@angular/router';
import { StakeholdersService } from 'src/app/infrastructure/stakeholders.service';

@Component({
  selector: 'xp-tour-list',
  templateUrl: './tour-list.component.html',
  styleUrls: ['./tour-list.component.css']
})
export class TourListComponent implements OnInit, OnDestroy {
  myTours$: Observable<Tour[]> | null = null;
  publishedTours$: Observable<Tour[]>;
  isCreator = false;
  private userSubscription?: Subscription;
  readonly defaultAvatar = 'assets/images/default_profile.png';
  private authorCache = new Map<number, AuthorMeta>();
  private pendingAuthorRequests = new Set<number>();
  private toursByAuthor = new Map<number, Tour[]>();

  constructor(
    private tourService: TourService,
    private authService: AuthService,
    private router: Router,
    private stakeholdersService: StakeholdersService
  ) {}

  ngOnInit(): void {
    this.publishedTours$ = this.decorateToursStream(this.tourService.getPublishedTours());

    this.userSubscription = this.authService.user$.subscribe(user => {
      const role = (user.role || '').toLowerCase();
      this.isCreator = role === 'author' || role === 'guide';
      this.myTours$ = this.isCreator ? this.decorateToursStream(this.tourService.getAuthorTours()) : null;
    });
  }

  ngOnDestroy(): void {
    this.userSubscription?.unsubscribe();
  }

  viewTour(tour: Tour): void {
    if (!tour?.id) {
      return;
    }
    this.router.navigate(['/tours', tour.id]);
  }

  private decorateToursStream(stream: Observable<Tour[]>): Observable<Tour[]> {
    return stream.pipe(
      tap(tours => tours.forEach(tour => this.hydrateAuthorMeta(tour)))
    );
  }

  private hydrateAuthorMeta(tour: Tour): void {
    if (tour?.authorId === null || tour?.authorId === undefined) {
      return;
    }

    if (tour.authorId <= 0) {
      tour.authorProfileImage = this.defaultAvatar;
      tour.authorName = 'Unknown author';
      return;
    }

    if (!this.toursByAuthor.has(tour.authorId)) {
      this.toursByAuthor.set(tour.authorId, []);
    }
    const trackedTours = this.toursByAuthor.get(tour.authorId)!;
    if (!trackedTours.includes(tour)) {
      trackedTours.push(tour);
    }

    const cached = this.authorCache.get(tour.authorId);
    if (cached) {
      this.assignAuthorMeta(tour.authorId, cached);
      return;
    }

    tour.authorProfileImage = this.defaultAvatar;

    if (this.pendingAuthorRequests.has(tour.authorId)) {
      return;
    }

    this.pendingAuthorRequests.add(tour.authorId);

    this.stakeholdersService.getUserById(tour.authorId).subscribe({
      next: (user) => {
        const meta: AuthorMeta = {
          name: `${user.first_name || ''} ${user.last_name || ''}`.trim(),
          username: user.username,
          avatar: user.profile_image || this.defaultAvatar
        };
        this.authorCache.set(tour.authorId, meta);
        this.assignAuthorMeta(tour.authorId, meta);
        this.pendingAuthorRequests.delete(tour.authorId);
      },
      error: () => {
        const fallback: AuthorMeta = { avatar: this.defaultAvatar };
        this.authorCache.set(tour.authorId, fallback);
        this.assignAuthorMeta(tour.authorId, fallback);
        this.pendingAuthorRequests.delete(tour.authorId);
      }
    });
  }

  private assignAuthorMeta(authorId: number, meta: AuthorMeta): void {
    const tours = this.toursByAuthor.get(authorId) || [];
    tours.forEach(t => {
      t.authorProfileImage = meta.avatar;
      if (meta.name) {
        t.authorName = meta.name;
      }
      if (meta.username) {
        t.authorUsername = meta.username;
      }
    });
  }
}

interface AuthorMeta {
  name?: string;
  username?: string;
  avatar: string;
}