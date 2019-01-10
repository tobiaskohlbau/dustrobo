import { Routes } from '@angular/router';
import { DashboardComponent } from './dashboard/dashboard.component';
import { MusicComponent } from './music/music.component';

export const routes: Routes = [
  { path: 'dashboard', component: DashboardComponent },
  { path: 'music', component: MusicComponent },
  { path: '', redirectTo: '/dashboard', pathMatch: 'full' },
  { path: '**', redirectTo: '/dashboard' }
];