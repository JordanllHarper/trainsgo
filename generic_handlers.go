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
) (HttpResponse, HttpError) {
	query := req.URL.Query()
	hasId := query.Has("id")

	if hasId {
		parsedId, err := uuid.Parse(query.Get("id"))
		if err != nil {
			return nil, badId(parsedId.String())
		}
		value, storeErr := store.GetById(parsedId)

		if storeErr != nil {
			return nil, storeErr
		}
		return statusOK{value}, nil
	} else {
		all, storeErr := store.All()
		if storeErr != nil {
			return nil, internalServerError{storeErr}
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
		return nil, internalServerError{delErr}
	}

	return statusOK{nil}, nil

}
