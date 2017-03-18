import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { HttpModule } from '@angular/http';

import { AppComponent } from './app.component';
import { VenueService } from './venue.service';
import { MapModule } from './map/map.module';

@NgModule({
  declarations: [
    AppComponent
  ],
  imports: [
    BrowserModule,
    FormsModule,
    HttpModule,
    MapModule
  ],
  providers: [VenueService],
  bootstrap: [AppComponent]
})
export class AppModule { }
