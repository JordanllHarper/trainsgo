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
	err HttpError,
) {
	if err != nil {
		http.Error(w, err.Error(), err.HttpCode())
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
