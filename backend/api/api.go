package api

import (
  "net/http"
)

type travelTimesResponse struct {
	stopTravelTimeMap map[int]string
}

// takes stop id from request, uses DAO to get traveltimes from postgres
func travelTimes(w http.ResponseWriter, r *http.Request) {

}

