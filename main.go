package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

	trainStore := newTrainStoreLocal()
	stationStore := newStationStoreLocal()

	trainHandler := newTrainHandler(trainStore, stationStore)
	stationHandler := newStationHandler(stationStore)

	st1 := newStation(newPosition(0, 0), "Station 1", 3)
	stationStore.register(st1)

	http.Handle("/trains", &trainHandler)
	http.Handle("/stations", &stationHandler)

	port := ":8080"
	fmt.Printf("Listening on: %s\n", port)
	log.Fatalln(http.ListenAndServe(port, nil))

}
