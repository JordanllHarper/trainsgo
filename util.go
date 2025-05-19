package main

import "github.com/google/uuid"

func getByStringId[V any](values StoreReader[V], id string) (v V, jsonErr any) {
	stId, err := uuid.Parse(id)
	if err != nil {
		return v, errBadId(id)
	}

	v, stErr := values.GetById(stId)
	if stErr != nil {
		return v, errIdDoesntExist(stId)
	}
	return v, nil
}
