import { Injectable } from '@angular/core';
import { Http } from '@angular/http';
import { Venue } from './venue';
import { VenueReview } from './venue-review';
import 'rxjs/add/operator/toPromise';

@Injectable()
export class VenueService {

  constructor(private http: Http) { }

  getVenue(googlePlacesId: string): Promise<Venue> {
    return this.http
        .get('/api/venue/' + googlePlacesId)
        .toPromise()
        .then(response => response.json().data as Venue);
  }

  getVenues(googlePlacesIds: string[]): Promise<Venue[]> {
    return this.http
        .get('/api/venue/list?venues=' + googlePlacesIds.join(','))
        .toPromise()
        .then(response => response.json().data as Venue[]);
  }

  postVenueReview(review: VenueReview): Promise<Venue> {
    return this.http
        .post('/api/venue/' + review.GoogleMapsId, review)
        .toPromise()
        .then(response => response.json().data as Venue);
  }

}
