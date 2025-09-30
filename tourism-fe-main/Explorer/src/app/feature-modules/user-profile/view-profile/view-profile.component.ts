import { Component } from '@angular/core';
import { User } from '../profile/model/profile.model';
import { ActivatedRoute } from '@angular/router';
import { StakeholdersService } from 'src/app/infrastructure/stakeholders.service';

@Component({
  selector: 'xp-view-profile',
  templateUrl: './view-profile.component.html',
  styleUrls: ['./view-profile.component.css']
})
export class ViewProfileComponent {
  user: User;
  isLoading = true;

  constructor(
    private route: ActivatedRoute,
    private stakeholdersService: StakeholdersService
  ) { }

  ngOnInit(): void {
    this.isLoading = true;
    // Nema 'if/else', uvek čitamo ID
    this.route.params.subscribe(params => {
      const userId = params['id'];
      this.stakeholdersService.getUserById(userId).subscribe(user => {
        this.user = user;
        this.isLoading = false;
      });
    });
  }
}
