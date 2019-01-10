import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from 'src/environments/environment';

@Injectable({
  providedIn: 'root'
})
export class MusicService {
  private readonly API_URL = environment.apiURL + "/music";

  constructor(private httpClient: HttpClient) { }

  play(url: string): Observable<any> {
    return this.httpClient.post(this.API_URL+"/play", {url: url});
  }

  pause(): Observable<any> {
    return this.httpClient.get(this.API_URL+"/pause");
  }

  stop(): Observable<any> {
    return this.httpClient.get(this.API_URL+"/stop");
  }

  volume(volume: number): Observable<any> {
    return this.httpClient.post(this.API_URL+"/volume", {volume: volume});
  }
}
