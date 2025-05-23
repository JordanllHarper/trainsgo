package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func serveJson(
	w http.ResponseWriter,
	method string,
	response HttpResponse,
	err error,
) {
	if err != nil {
		httpErr := mapToHttpErr(err)
		log.Printf("Received error %s request: %v\n", method, err)
		http.Error(w, httpErr.Error(), httpErr.HttpCode())
		return
	}

	log.Printf("Received %s request\n", method)

	w.WriteHeader(response.HttpCode())
	if response.Body() != nil {
		err := json.NewEncoder(w).Encode(response.Body())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
