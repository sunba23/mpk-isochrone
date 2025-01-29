package api

import (
	"net/http"
)

// API to take one stop from map, then run computing distances then return distances
func distances(w http.ResponseWriter, r *http.Request) {

}

func Run() {
  // TODO add middleware for stripping trailing slashes
  http.HandleFunc("/distances/", distances)
  http.ListenAndServe(":8080", nil)
}
