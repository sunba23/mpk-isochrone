package api

import (
	"log"
	"net/http"
)

func Run() {
	// wrap travelTimesAndRoute in middleware(s)
	timeRouteHandler := restrictMethod(http.MethodGet, stripSlashes(http.HandlerFunc(travelData)))
	// register the handler(s)
	http.Handle("/traveldata", timeRouteHandler)
	// run server
	log.Fatal(http.ListenAndServe(":8080", nil))
}
