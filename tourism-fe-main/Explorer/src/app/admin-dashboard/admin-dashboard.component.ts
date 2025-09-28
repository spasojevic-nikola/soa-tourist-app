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
  successMessage = '';

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

  blockUser(userId: number): void {
    const token = this.tokenStorage.getAccessToken();
    
    if (!token) {
      this.errorMessage = 'Token not found. Please log in.';
      return;
    }

    const headers = { 'Authorization': `Bearer ${token}` };
    
    // Pozivamo auth-service endpoint za blokiranje
    this.http.post(`${environment.authApiHost}admin/users/${userId}/block`, {}, { headers }).subscribe({
      next: () => {
        this.successMessage = 'User successfully blocked.';
        this.errorMessage = '';
        // Refresh the user list to show updated status
        this.getAllUsers();
        // Clear success message after 3 seconds
        setTimeout(() => {
          this.successMessage = '';
        }, 3000);
      },
      error: (error) => {
        console.error('Failed to block user', error);
        if (error.status === 403) {
          this.errorMessage = 'Access Denied. You do not have administrator privileges.';
        } else if (error.status === 400) {
          this.errorMessage = 'Cannot block administrator.';
        } else {
          this.errorMessage = 'Failed to block user.';
        }
        this.successMessage = '';
      }
    });
  }

  unblockUser(userId: number): void {
    const token = this.tokenStorage.getAccessToken();
    
    if (!token) {
      this.errorMessage = 'Token not found. Please log in.';
      return;
    }

    const headers = { 'Authorization': `Bearer ${token}` };
    
    // Pozivamo auth-service endpoint za odblokiranje
    this.http.post(`${environment.authApiHost}admin/users/${userId}/unblock`, {}, { headers }).subscribe({
      next: () => {
        this.successMessage = 'User successfully unblocked.';
        this.errorMessage = '';
        // Refresh the user list to show updated status
        this.getAllUsers();
        // Clear success message after 3 seconds
        setTimeout(() => {
          this.successMessage = '';
        }, 3000);
      },
      error: (error) => {
        console.error('Failed to unblock user', error);
        if (error.status === 403) {
          this.errorMessage = 'Access Denied. You do not have administrator privileges.';
        } else {
          this.errorMessage = 'Failed to unblock user.';
        }
        this.successMessage = '';
      }
    });
  }
}