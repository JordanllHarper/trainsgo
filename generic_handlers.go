package main

import (
	"maps"
	"net/http"
	"net/url"
	"slices"

	"github.com/google/uuid"
)

const (
	// Supported query parameters
	id   = "id"
	name = "name"
)

func handleGetNameIdAll[V any](qry url.Values, s StoreIdNameable[V]) (HttpResponse, error) {
	if qry.Has(name) {
		return handleGetByName(qry.Get(name), s)
	}
	return handleGetIdAll(qry, s)
}

func handleGetIdAll[V any](qry url.Values, s StoreIDable[V]) (HttpResponse, error) {
	if qry.Has(id) {
		return handleGetById(qry.Get(id), s)
	}
	return handleGetAll(s)
}

func handleGetAll[V any](store Store[V]) (HttpResponse, error) {
	all, err := store.All()
	if err != nil {
		return nil, internalError{err}
	}
	storeValues := maps.Values(all)
	arr := slices.Collect(storeValues)
	return statusOKMult[V]{Values: arr}, nil
}

func handleGetById[V any](
	id string,
	store IDable[V],
) (HttpResponse, error) {
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

func handleGetByName[V any](
	name string,
	store Nameable[V],
) (HttpResponse, error) {
	value, err := store.GetByName(name)
	if err != nil {
		return nil, err
	}
	return statusOK{value}, nil
}

func handleDelete(req *http.Request, store Deleter) (HttpResponse, HttpError) {
	var t deleteBody
	if err := jsonDecode(req, &t); err != nil {
		return nil, malformedBody{}
	}
	id, err := uuid.Parse(t.Id)
	if err != nil {
		return nil, badId(t.Id)
	}
	if err = store.Delete(id); err != nil {
		return nil, internalError{err}
	}
	return statusOK{}, nil

}
