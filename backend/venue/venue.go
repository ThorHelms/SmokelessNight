package venue

import (
	"encoding/json"
	"fmt"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type DecliningValuation struct {
	Value   float64
	Changed time.Time
}

type ConstrictedDecliningValuation struct {
	DecliningValuation
	MaxValue float64
}

const (
	Unknown               = iota
	SmokingAllowed        = iota
	SmokingProhibited     = iota
	PartialSmokingAllowed = iota
)

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

const venue = "Venue"
const invalidKeyMsg = "Invalid key: "
const methodNotAllowedMsg = "Method not allowed: "

var isValidKeyRegexp = regexp.MustCompile("^[a-zA-Z0-9-]+$")

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
	} else if strings.EqualFold(r.Method, "patch") {
		if isValidKey(relUrl) {
			postVenueReviewHandler(relUrl, w, r)
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
	e := new(Venue)
	if err := datastore.Get(ctx, k, e); err != nil {
		if err == datastore.ErrNoSuchEntity {
			http.Error(w, err.Error(), 404)
		} else {
			http.Error(w, err.Error(), 500)
		}
		return
	}
	encoder := json.NewEncoder(w)
	encoder.Encode(e)
}

func listVenuesHandler(venueKeys []string, w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	venues := []*Venue{}
	for i := 0; i < len(venueKeys); i++ {
		if !isValidKey(venueKeys[i]) {
			http.Error(w, invalidKeyMsg+venueKeys[i], 400)
			return
		}
		k := datastore.NewKey(ctx, venue, venueKeys[i], 0, nil)
		e := new(Venue)
		if err := datastore.Get(ctx, k, e); err != nil {
			if err != datastore.ErrNoSuchEntity {
				http.Error(w, err.Error(), 500)
				return
			}
		} else {
			venues = append(venues, e)
		}
	}
	encoder := json.NewEncoder(w)
	encoder.Encode(venues)
}

func postVenueReviewHandler(relUrl string, w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Post venue review please: "+relUrl)
}
