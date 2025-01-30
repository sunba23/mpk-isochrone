package api

import (
	"fmt"
	"log"
	"net/http"
)

func Run() {
	// wrap travelTimesAndRoute in middleware(s)
	timeRouteHandler := restrictMethod(http.MethodGet, stripSlashes(http.HandlerFunc(travelData)))
	// register the handler(s)
	http.Handle("/traveldata", timeRouteHandler)
	// run server
  fmt.Println("Starting go API server.")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
