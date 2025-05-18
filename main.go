package main

import (
	"fmt"
	"log"
	"net/http"
)

type HandlerConfiguration struct {
	trainHandler   http.Handler
	stationHandler http.Handler
	lineHandler    http.Handler
}

func createDummyConfiguration() HandlerConfiguration {

	trainStore := newTrainStoreLocal()
	stationStore := newStationStoreLocal()
	lineStore := newLineStoreLocal()

	st1 := newStation(newPosition(0, 0), "Station 1", 3)
	st2 := newStation(newPosition(10, 10), "Station 2", 5)
	t1 := newTrain("Train 1", st1)

	trainStore.register(t1)
	stationStore.register(st1)
	stationStore.register(st2)
	lineStore.register(newLine(st1, st2, "Line 1"))

	trainHandlerLocal := trainHandlerLocal{trainStore, stationStore}
	lineHandlerLocal := lineHandlerLocal{lineStore, stationStore}

	return HandlerConfiguration{
		trainHandlerLocal,
		stationStore,
		lineHandlerLocal,
	}
}

func main() {
	config := createDummyConfiguration()
	{

		http.Handle("/trains", config.trainHandler)
		http.Handle("/stations", config.stationHandler)
		http.Handle("/line", config.lineHandler)
		// http.HandleFunc("/trip", func(w http.ResponseWriter, req *http.Request) {
		// 	handleTrip(w, req, tripStore)
		// })
	}
	{
		port := ":8080"
		fmt.Printf("Listening on: %s\n", port)
		log.Fatalln(http.ListenAndServe(port, nil))
	}
}
