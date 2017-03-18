package venue

import (
	"encoding/json"
	"github.com/gorilla/schema"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"math"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const (
	minAirQuality         = 0.0
	maxAirQuality         = 1.0
	depreciationValue     = 0.99
	day                   = time.Duration(24) * time.Hour
	Unknown               = 0
	SmokingAllowed        = 1
	SmokingProhibited     = 2
	PartialSmokingAllowed = 3
	venue                 = "Venue"
	invalidKeyMsg         = "Invalid key: "
	methodNotAllowedMsg   = "Method not allowed: "
	reviewNotSpecifiedMsg = "No review is specified"
	reviewNotInBoundMsg   = "Review-values not in bound"
)

var isValidKeyRegexp = regexp.MustCompile("^[a-zA-Z0-9-]+$")
var decoder = schema.NewDecoder()

type DecliningValuation struct {
	Value   float64
	Changed time.Time
}

type ConstrictedDecliningValuation struct {
	DecliningValuation
	MaxValue float64
}

type Venue struct {
	GoogleMapsId         string
	AirQualitySmokers    float64 `datastore:",noindex"`
	AirQualityNonsmokers float64 `datastore:",noindex"`
	SmokingPolicy        int     `datastore:",noindex"`
}

type VenueDataModel struct {
	Venue
	SmokingAllowed                  DecliningValuation            `datastore:",noindex"`
	SmokingProhibited               DecliningValuation            `datastore:",noindex"`
	PartialSmokingAllowed           DecliningValuation            `datastore:",noindex"`
	AirQualityCalculationSmokers    ConstrictedDecliningValuation `datastore:",noindex"`
	AirQualityCalculationNonsmokers ConstrictedDecliningValuation `datastore:",noindex"`
}

type VenueReview struct {
	Venue
	SmokingPolicySpecified        bool
	AirQualitySmokersSpecified    bool
	AirQualityNonsmokersSpecified bool
}

func daysSince(then time.Time, now time.Time) float64 {
	then = then.Truncate(day)
	now = now.Truncate(day)
	diff := now.Sub(then)
	return math.Floor(diff.Hours() / 24.0)
}

func (dv *DecliningValuation) Depreciate(now time.Time) {
	days := daysSince(dv.Changed, now)
	if days < 0.5 {
		return
	}
	dv.Value = math.Pow(depreciationValue, days) * dv.Value
	dv.Changed = now
}

func (cdv *ConstrictedDecliningValuation) DepreciateAll(now time.Time) {
	days := daysSince(cdv.Changed, now)
	if days < 0.5 {
		return
	}
	cdv.Value = math.Pow(depreciationValue, days) * cdv.Value
	cdv.MaxValue = math.Pow(depreciationValue, days) * cdv.MaxValue
	cdv.Changed = now
}

func (cdv ConstrictedDecliningValuation) WeightedValue() float64 {
	if cdv.MaxValue == 0.0 {
		return -1.0
	}
	return cdv.Value / cdv.MaxValue
}

func (venue *VenueDataModel) Reset() {
	now := time.Now()
	venue.AirQualityCalculationNonsmokers.Changed = now
	venue.AirQualityCalculationSmokers.Changed = now
	venue.PartialSmokingAllowed.Changed = now
	venue.SmokingAllowed.Changed = now
	venue.SmokingProhibited.Changed = now
	venue.SmokingPolicy = Unknown
}

func (venue *VenueDataModel) Depreciate() {
	now := time.Now()
	venue.AirQualityCalculationNonsmokers.DepreciateAll(now)
	venue.AirQualityCalculationSmokers.DepreciateAll(now)
	venue.PartialSmokingAllowed.Depreciate(now)
	venue.SmokingAllowed.Depreciate(now)
	venue.SmokingProhibited.Depreciate(now)
}

func (venue *VenueDataModel) AddReview(review *VenueReview) {
	if review.AirQualityNonsmokersSpecified {
		venue.AirQualityCalculationNonsmokers.Value += review.AirQualityNonsmokers
		venue.AirQualityCalculationNonsmokers.MaxValue += maxAirQuality
		venue.AirQualityNonsmokers = venue.AirQualityCalculationNonsmokers.WeightedValue()
	}
	if review.AirQualitySmokersSpecified {
		venue.AirQualityCalculationSmokers.Value += review.AirQualitySmokers
		venue.AirQualityCalculationSmokers.MaxValue += maxAirQuality
		venue.AirQualitySmokers = venue.AirQualityCalculationSmokers.WeightedValue()
	}
	if review.SmokingPolicySpecified {
		switch review.SmokingPolicy {
		case SmokingAllowed:
			venue.SmokingAllowed.Value += 1.0
		case SmokingProhibited:
			venue.SmokingProhibited.Value += 1.0
		case PartialSmokingAllowed:
			venue.PartialSmokingAllowed.Value += 1.0
		}
		if venue.SmokingAllowed.Value > venue.SmokingProhibited.Value && venue.SmokingAllowed.Value > venue.PartialSmokingAllowed.Value {
			venue.SmokingPolicy = SmokingAllowed
		} else if venue.SmokingProhibited.Value > venue.SmokingAllowed.Value && venue.SmokingProhibited.Value > venue.PartialSmokingAllowed.Value {
			venue.SmokingPolicy = SmokingProhibited
		} else if venue.PartialSmokingAllowed.Value > venue.SmokingAllowed.Value && venue.PartialSmokingAllowed.Value > venue.SmokingProhibited.Value {
			venue.SmokingPolicy = PartialSmokingAllowed
		} else {
			venue.SmokingPolicy = Unknown
		}
	}
}

func (review VenueReview) Specified() bool {
	return review.AirQualityNonsmokersSpecified || review.AirQualitySmokersSpecified || review.SmokingPolicySpecified
}

func (review VenueReview) ValuesInBounds() bool {
	if review.AirQualityNonsmokersSpecified && (review.AirQualityNonsmokers < minAirQuality || review.AirQualityNonsmokers > maxAirQuality) {
		return false
	}
	if review.AirQualitySmokersSpecified && (review.AirQualitySmokers < minAirQuality || review.AirQualitySmokers > maxAirQuality) {
		return false
	}
	if review.SmokingPolicySpecified && review.SmokingPolicy != SmokingAllowed && review.SmokingPolicy != SmokingProhibited && review.SmokingPolicy != PartialSmokingAllowed {
		return false
	}
	return true
}

func isValidKey(relUrl string) bool {
	return isValidKeyRegexp.MatchString(relUrl)
}

func Handler(relUrl string, w http.ResponseWriter, r *http.Request) {
	if strings.EqualFold(r.Method, "get") {
		if strings.EqualFold(relUrl, "list") {
			venuesRaw := r.URL.Query().Get("venues")
			venues := strings.Split(venuesRaw, ",")
			listVenuesHandler(venues, w, r)
		} else if isValidKey(relUrl) {
			getVenueHandler(relUrl, w, r)
		} else {
			http.Error(w, invalidKeyMsg+relUrl, 400)
		}
	} else if strings.EqualFold(r.Method, "post") {
		if isValidKey(relUrl) {
			if err := r.ParseForm(); err != nil {
				http.Error(w, err.Error(), 400)
				return
			}
			review := new(VenueReview)
			if err := decoder.Decode(review, r.PostForm); err != nil {
				http.Error(w, err.Error(), 400)
				return
			}
			if !review.Specified() {
				http.Error(w, reviewNotSpecifiedMsg, 400)
				return
			}
			if !review.ValuesInBounds() {
				http.Error(w, reviewNotInBoundMsg, 400)
				return
			}
			postVenueReviewHandler(relUrl, review, w, r)
		} else {
			http.Error(w, invalidKeyMsg+relUrl, 400)
		}
	} else {
		http.Error(w, methodNotAllowedMsg+r.Method, 405)
	}
}

func getVenueHandler(relUrl string, w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	k := datastore.NewKey(ctx, venue, relUrl, 0, nil)
	e := new(VenueDataModel)
	if err := datastore.Get(ctx, k, e); err != nil {
		if err == datastore.ErrNoSuchEntity {
			http.Error(w, err.Error(), 404)
		} else {
			http.Error(w, err.Error(), 500)
		}
		return
	}
	encoder := json.NewEncoder(w)
	encoder.Encode(e.Venue)
}

func listVenuesHandler(venueKeys []string, w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	venues := []Venue{}
	for i := 0; i < len(venueKeys); i++ {
		if !isValidKey(venueKeys[i]) {
			http.Error(w, invalidKeyMsg+venueKeys[i], 400)
			return
		}
		k := datastore.NewKey(ctx, venue, venueKeys[i], 0, nil)
		e := new(VenueDataModel)
		if err := datastore.Get(ctx, k, e); err != nil {
			if err != datastore.ErrNoSuchEntity {
				http.Error(w, err.Error(), 500)
				return
			}
		} else {
			venues = append(venues, e.Venue)
		}
	}
	encoder := json.NewEncoder(w)
	encoder.Encode(venues)
}

func postVenueReviewHandler(relUrl string, review *VenueReview, w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	k := datastore.NewKey(ctx, venue, relUrl, 0, nil)
	e := new(VenueDataModel)
	if err := datastore.Get(ctx, k, e); err != nil {
		if err == datastore.ErrNoSuchEntity {
			e.Reset()
			e.GoogleMapsId = relUrl
		} else {
			http.Error(w, err.Error(), 500)
			return
		}
	}
	e.Depreciate()
	e.AddReview(review)
	if _, err := datastore.Put(ctx, k, e); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	encoder := json.NewEncoder(w)
	encoder.Encode(e.Venue)
}
