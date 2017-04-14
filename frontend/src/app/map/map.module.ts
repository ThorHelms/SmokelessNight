import { BrowserModule } from '@angular/platform-browser';
import { NgModule, ApplicationRef } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { MapComponent } from './map.component';

import { AgmCoreModule } from '@agm/core'

@NgModule({
  imports: [
    BrowserModule,
    CommonModule,
    FormsModule,
    AgmCoreModule.forRoot({
      apiKey: 'AIzaSyBcBU0pxPwNaSiWcdttNm1YwaN3Dl0SjEk'
    })
  ],
  providers: [],
  declarations: [ MapComponent ],
  exports: [ MapComponent ]
})
export class MapModule {}
