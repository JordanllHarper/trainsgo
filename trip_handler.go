package main

import "net/http"

func handleTrip(w http.ResponseWriter, req *http.Request, trc tripReaderCoordinator) {
	method := req.Method

	switch method {
	case "POST":
		handleTripPost(w, req, trc)
	case "DELETE":
		handleTripDelete(w, req, trc)
	default:
		http.Error(w, http.ErrNotSupported.Error(), http.StatusMethodNotAllowed)
	}
}

func handleTripDelete(w http.ResponseWriter, req *http.Request, trc tripReaderCoordinator) {
	// TODO: Cancel a trip if possible
	panic("unimplemented")
}

func handleTripPost(w http.ResponseWriter, req *http.Request, trc tripCoordinator) {

	// TODO: Create a new trip
	// _ := trc.scheduleTrip()
	panic("unimplemented")
}

func handleTripPut(w http.ResponseWriter, req *http.Request, trc tripCoordinator) {

}
