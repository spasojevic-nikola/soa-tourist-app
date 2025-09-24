import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { HomeComponent } from 'src/app/feature-modules/layout/home/home.component';
import { LoginComponent } from '../auth/login/login.component';
import { AuthGuard } from '../auth/auth.guard';
import { RegistrationComponent } from '../auth/registration/registration.component';
import { ProfileComponent } from 'src/app/feature-modules/user-profile/profile/profile.component';
import { AdminDashboardComponent } from 'src/app/admin-dashboard/admin-dashboard.component';

const routes: Routes = [
  {path: 'home', component: HomeComponent},
  {path: 'login', component: LoginComponent},
  {path: 'register', component: RegistrationComponent},
  {
    path: 'profile', 
    loadChildren: () => import('../../feature-modules/user-profile/user-profile.module').then(m => m.UserProfileModule),
    canActivate: [AuthGuard] // DODAJTE OVO
  },
  { path: 'admin/users', component: AdminDashboardComponent },  
  {
    path: 'blog',
    loadChildren: () => import('../../feature-modules/blog/blog.module').then(m => m.BlogModule),
    canActivate: [AuthGuard]
  }

];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
