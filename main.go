package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	trainStore := newTrainStoreLocal()
	stationStore := newStationStoreLocal()
	http.Handle("/trains", &trainHandler{trainStore, stationStore})
	http.Handle("/stations", &stationHandler{stationStore})

	port := ":8080"
	fmt.Printf("Listening on: %s\n", port)
	log.Fatalln(http.ListenAndServe(port, nil))

}
