package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func serveJson(
	w http.ResponseWriter,
	method string,
	code int,
	body any,
) {
	log.Printf("Received %s request\n", method)

	w.WriteHeader(code)
	if body != nil {
		err := json.NewEncoder(w).Encode(body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
