package main

import (
	"github.com/google/uuid"
)

func getByStringId[V any](values StoreIDable[V], id string) (V, error) {
	var v V
	stId, err := uuid.Parse(id)
	if err != nil {
		return v, badId(id)
	}

	v, stErr := values.GetById(stId)
	if stErr != nil {
		return v, err
	}
	return v, nil
}

func sliceUnorderedRemove[T any](s []T, i int) []T {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
