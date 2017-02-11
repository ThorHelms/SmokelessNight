package api

import (
	"net/http"
	"strings"
	"venue"
)

func handleFunc(url string, handler func(string, http.ResponseWriter, *http.Request)) {
	http.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		relUrl := strings.Replace(r.URL.Path, url, "", 1)
		handler(relUrl, w, r)
	})
}

func init() {
	handleFunc("/api/venue/", venue.Handler)
}
