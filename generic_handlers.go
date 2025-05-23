package main

import (
	"encoding/json"
	"maps"
	"net/http"
	"slices"

	"github.com/google/uuid"
)

func handleGet[V any](
	req *http.Request,
	store StoreReader[V],
) (HttpResponse, error) {
	query := req.URL.Query()

	if query.Has("id") {
		id := query.Get("id")
		parsedId, err := uuid.Parse(id)
		if err != nil {
			return nil, badId(id)
		}

		value, err := store.GetById(parsedId)
		if err != nil {
			return nil, err
		}
		return statusOK{value}, nil
	}

	all, err := store.All()
	if err != nil {
		return nil, internalError{err}
	}

	storeValues := maps.Values(all)
	arr := slices.Collect(storeValues)

	response := struct {
		Values []V `json:"values"`
	}{
		Values: arr,
	}
	return statusOK{response}, nil
}

func handleDelete(req *http.Request, store StoreDeleter) (HttpResponse, HttpError) {
	var t deleteBody

	if err := json.NewDecoder(req.Body).Decode(&t); err != nil {
		return nil, malformedBody{}
	}

	id, err := uuid.Parse(t.Id)
	if err != nil {
		return nil, badId(t.Id)
	}

	if delErr := store.Delete(id); delErr != nil {
		return nil, internalError{delErr}
	}

	return statusOK{nil}, nil

}
