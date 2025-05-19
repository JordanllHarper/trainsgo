package main

import (
	"fmt"
	"log"
	"net/http"
)

type HandlerConfiguration struct {
	trainHandler   http.HandlerFunc
	stationHandler http.HandlerFunc
	lineHandler    http.HandlerFunc
}

func createDummyConfiguration() HandlerConfiguration {

	trainStore := NewTrainStoreLocal()
	stationStore := NewStationStoreLocal()
	lineStore := NewLineStoreLocal()

	st1 := NewStation(NewPosition(0, 0), "Station 1", 3)
	st2 := NewStation(NewPosition(10, 10), "Station 2", 5)
	t1 := NewTrain("Train 1", st1)

	trainStore.register(t1)
	stationStore.register(st1)
	stationStore.register(st2)
	lineStore.register(newLine(st1, st2, "Line 1"))

	trainHandlerLocal := trainHandlerLocal{trainStore, stationStore}
	lineHandlerLocal := lineHandlerLocal{lineStore, stationStore}

	return HandlerConfiguration{
		stationHandler: stationStore.ServeHTTP,
		trainHandler:   trainHandlerLocal.ServeHTTP,
		lineHandler:    lineHandlerLocal.ServeHTTP,
	}
}

func main() {
	config := createDummyConfiguration()
	{

		http.HandleFunc("/trains", config.trainHandler)
		http.HandleFunc("/stations", config.stationHandler)
		http.HandleFunc("/line", config.lineHandler)
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
