import { Component, OnInit } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from 'src/env/environment';
import { TokenStorage } from 'src/app/infrastructure/auth/jwt/token.service'; 

@Component({
  selector: 'app-admin-dashboard',
  templateUrl: './admin-dashboard.component.html',
  styleUrls: ['./admin-dashboard.component.css']
})
export class AdminDashboardComponent implements OnInit {
  users: any[] = [];
  errorMessage = '';

  constructor(private http: HttpClient, private tokenStorage: TokenStorage) { }

  ngOnInit(): void {
    this.getAllUsers();
  }

  getAllUsers(): void {
  const token = this.tokenStorage.getAccessToken();

    if (!token) {
      this.errorMessage = 'Token not found. Please log in.';
      return;
    }

    const headers = { 'Authorization': `Bearer ${token}` };
    
    // Pozivamo stakeholders-service
    this.http.get<any[]>(`${environment.stakeholdersApiHost}/admin/users`, { headers }).subscribe({
      next: (data) => {
        this.users = data;
        this.errorMessage = '';
      },
      error: (error) => {
        console.error('Failed to fetch users', error);
        if (error.status === 403) {
          this.errorMessage = 'Access Denied. You do not have administrator privileges.';
        } else {
          this.errorMessage = 'Failed to load user data.';
        }
      }
    });
  }
}