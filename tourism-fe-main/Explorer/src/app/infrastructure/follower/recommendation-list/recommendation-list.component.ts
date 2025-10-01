import { Component } from '@angular/core';
import { Recommendation } from 'src/app/shared/model/recommendation.model';
import { FollowerService } from '../follower.service';
import { catchError, EMPTY, finalize } from 'rxjs';

@Component({
  selector: 'xp-recommendation-list',
  templateUrl: './recommendation-list.component.html',
  styleUrls: ['./recommendation-list.component.css']
})
export class RecommendationListComponent {

  recommendations: Recommendation[] = [];
  isLoading = true;
  error: string | null = null;

  constructor(private followerService: FollowerService) { }

  ngOnInit(): void {
    this.loadRecommendations();
  }
  loadRecommendations(): void {
    this.isLoading = true;
    this.error = null;

    //  Ovde bi u realnoj aplikaciji nakon dobijanja ID-jeva (recommendations)
    // sledio jos jedan HTTP poziv User Service-u da se dohvate username, slika, itd.
    // Za sada, samo dohvatamo ID-jeve i score.
    this.followerService.getRecommendations()
      .pipe(
        finalize(() => this.isLoading = false),
        catchError(err => {
          console.error("Greška pri dohvatanju preporuka:", err);
          this.error = 'Greška pri dohvatanju preporuka. Pokušajte ponovo.';
          return EMPTY; // Zaustavlja tok
        })
      )
      .subscribe(data => {
        this.recommendations = data;
      });
  }

  followUser(userId: number): void {
    this.followerService.follow(userId).subscribe({
      next: () => {
        alert(`Uspešno zapraćen korisnik ID: ${userId}`);
        // Uklonite korisnika iz preporuka ili osvežite listu
        this.recommendations = this.recommendations.filter(r => r.userId !== userId);
      },
      error: (err) => {
        console.error("Greška pri praćenju:", err);
        alert('Došlo je do greške prilikom praćenja.');
      }
    });
  }
}
