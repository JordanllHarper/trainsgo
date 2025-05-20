package main

import "github.com/google/uuid"

func getByStringId[V any](values StoreReader[V], id string) (V, HttpError) {
	var v V
	stId, err := uuid.Parse(id)
	if err != nil {
		return v, badId(id)
	}

	v, stErr := values.GetById(stId)
	if stErr != nil {
		return v, stErr
	}
	return v, nil
}
