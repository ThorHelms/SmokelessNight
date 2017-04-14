import { Component } from '@angular/core';
import { AgmCoreModule } from '@agm/core';

@Component({
  selector: 'app-map',
  templateUrl: './map.component.html',
  styleUrls: ['./map.component.scss']
})

export class MapComponent {
  title = 'My first angular2-google-maps project';
  lat = 55.6760968;
  lng = 12.5683371;
}
