import { Component, OnInit, ViewChild } from '@angular/core';
import { MusicService } from './music.service';
import { MatSlider } from '@angular/material';

@Component({
  selector: 'app-music',
  templateUrl: './music.component.html',
  styleUrls: ['./music.component.scss']
})
export class MusicComponent implements OnInit {
  volume: number;
  paused: boolean = false;

  constructor(private musicService: MusicService) { }

  ngOnInit() {
  }

  play(url: string) {
    this.musicService.play(url).subscribe(data => {}, error => console.log(error));
  }

  pause() {
    this.musicService.pause().subscribe(data => {this.paused = !this.paused;}, error => console.log(error));
  }

  stop() {
    this.musicService.stop().subscribe(data => {}, error => console.log(error));
  }

  onVolume(event) {
    this.musicService.volume(event.value).subscribe(data => {}, error => console.log(error));
  }

}
