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
	store storeReader[V],
) (int, any) {
	query := req.URL.Query()
	hasId := query.Has("id")

	if hasId {
		parsedId, err := uuid.Parse(query.Get("id"))
		if err != nil {
			return http.StatusInternalServerError, nil
		}
		value, storeErr := store.getById(parsedId)

		if storeErr != nil {
			switch storeErr.code {
			case StoreReaderErrIdNotFound:
				return http.StatusNotFound, nil
			default:
				return http.StatusInternalServerError, nil
			}
		}
		return http.StatusOK, value
	} else {
		all, storeErr := store.all()
		if storeErr != nil {
			return http.StatusInternalServerError, nil
		}

		storeValues := maps.Values(all)
		arr := slices.Collect(storeValues)

		response := struct {
			Values []V `json:"values"`
		}{
			Values: arr,
		}
		return http.StatusOK, response
	}
}

func handleDelete(req *http.Request, store storeDeleter) (int, any) {
	var t deleteBody

	if err := json.NewDecoder(req.Body).Decode(&t); err != nil {
		return http.StatusBadRequest, nil
	}

	id, err := uuid.Parse(t.Id)
	if err != nil {
		return http.StatusBadRequest, nil
	}

	if delErr := store.delete(id); delErr != nil {
		return http.StatusInternalServerError, nil
	}

	return http.StatusOK, nil

}
