package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

	trainStore := newTrainStoreLocal()
	stationStore := newStationStoreLocal()
	lineStore := newLineStoreLocal()

	trainHandler := newTrainHandler(trainStore, stationStore)
	stationHandler := newStationHandler(stationStore)
	navHandler := newNavHandler(lineStore, stationStore)

	st1 := newStation(newPosition(0, 0), "Station 1", 3)
	st2 := newStation(newPosition(10, 10), "Station 2", 5)
	stationStore.register(st1)
	stationStore.register(st2)
	fmt.Println(st1.E.Id, st2.E.Id)

	http.Handle("/trains", &trainHandler)
	http.Handle("/stations", &stationHandler)
	http.Handle("/nav", &navHandler)

	port := ":8080"
	fmt.Printf("Listening on: %s\n", port)
	log.Fatalln(http.ListenAndServe(port, nil))

}
