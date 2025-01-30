package api

import (
	"log"
	"net/http"
)

func Run() {
	// wrap travelTimes in middleware(s)
	travelHandler := restrictMethod(http.MethodGet, stripSlashes(http.HandlerFunc(travelTimes)))
  // register the handler(s)
	http.Handle("/traveltime/{stop_id}", travelHandler)
  // run server
	log.Fatal(http.ListenAndServe(":8080", nil))
}
