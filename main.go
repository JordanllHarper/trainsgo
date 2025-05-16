package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	trainStore := newTrainStoreLocal()
	stationStore := newStationStoreLocal()

	st1 := newStation(newPosition(0, 0), "Station 1", 3)
	stationStore.register(st1)

	http.Handle("/trains", &trainHandler{trainStore, stationStore})
	http.Handle("/stations", &stationHandler{stationStore})

	port := ":8080"
	fmt.Printf("Listening on: %s\n", port)
	log.Fatalln(http.ListenAndServe(port, nil))

}
