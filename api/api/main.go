package api

import (
	"fmt"
	"log"
	"net/http"
)

func Run() {
	// wrap functions in middlewares
	travelDataHandler := enableCORS(restrictMethod(http.MethodGet, stripSlashes(http.HandlerFunc(travelData))))
	stopsDetailsHandler := enableCORS(restrictMethod(http.MethodGet, stripSlashes(http.HandlerFunc(stopsDetails))))
	// register the handlers
	http.Handle("/traveldata", travelDataHandler)
	http.Handle("/stops/details", stopsDetailsHandler)
	// run server
	fmt.Println("Starting go API server.")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
