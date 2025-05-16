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

	{
		st1 := newStation(newPosition(0, 0), "Station 1", 3)
		st2 := newStation(newPosition(10, 10), "Station 2", 5)
		t1 := newTrain("Train 1", st1)

		// Registering a bunch of dummy data to test scheduling a train
		trainStore.register(t1)
		stationStore.register(st1)
		stationStore.register(st2)
		lineStore.register(newLine(st1, st2, "Line 1"))

		fmt.Println(st1.E.Id, st2.E.Id)

		setupHandlers(trainStore, stationStore, lineStore)
	}

	port := ":8080"
	fmt.Printf("Listening on: %s\n", port)
	log.Fatalln(http.ListenAndServe(port, nil))
}

func setupHandlers(
	trw storeReaderWriter[Train],
	srw storeReaderWriter[Station],
	lrw storeReaderWriter[Line],
) {

	http.HandleFunc("/trains", func(w http.ResponseWriter, req *http.Request) {
		serve(
			w,
			req,
			trw,
			func() { handleTrainPost(w, req, trw, srw) },
		)
	})
	http.HandleFunc("/stations", func(w http.ResponseWriter, req *http.Request) {
		serve(
			w,
			req,
			srw,
			func() { handleStationPost(w, req, srw) },
		)
	},
	)
	http.HandleFunc("/nav", func(w http.ResponseWriter, req *http.Request) {
		serve(
			w,
			req,
			lrw,
			func() { handleNavPost(w, req, lrw, srw) },
		)
	})
}
