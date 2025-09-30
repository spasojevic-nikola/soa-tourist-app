import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { StakeholdersService } from 'src/app/infrastructure/stakeholders.service';
import { User } from '../user-profile/profile/model/profile.model';

@Component({
  selector: 'app-search-results',
  templateUrl: './search-results.component.html',
  styleUrls: ['./search-results.component.css']
})
export class SearchResultsComponent implements OnInit {
  
  users: User[] = [];
  searchQuery: string = '';
  isLoading: boolean = true;
  noResults: boolean = false;

  constructor(
    private route: ActivatedRoute,
    private stakeholdersService: StakeholdersService
  ) { }

  ngOnInit(): void {
    this.route.queryParams.subscribe(params => {
      this.searchQuery = params['q'];
      if (this.searchQuery) {
        this.isLoading = true;
        this.noResults = false;
        this.users = [];

        this.stakeholdersService.searchUsers(this.searchQuery).subscribe({
          next: (foundUsers) => {
          //  console.log(foundUsers); 
            this.users = foundUsers;
            this.noResults = foundUsers.length === 0;
            this.isLoading = false;
          },
          error: (err) => {
            console.error('Error fetching search results:', err);
            this.isLoading = false;
            this.noResults = true;
          }
        });
      }
    });
  }
}