package main

import (
	"github.com/google/uuid"
)

type (
	/*
		Represents a path (a route) from a CurrentId to a DestId.

		NextId represents the next station or intersection or whatever to go to to continue the route.
	*/
	Route struct {
		RouteId   Id
		CurrentId Id
		DestId    Id
		NextId    Id
	}
)

func NewRoute(
	srcId Id,
	destId Id,
	nextId Id,
) Route {
	return Route{
		uuid.New(),
		srcId,
		destId,
		nextId,
	}
}

//

type (
	RouteStore interface {
		StoreReaderDeleter[Route]
		GetByCurrDest(current, dest Id) (Route, StoreError)
	}

	routeStoreLocal map[Id]Route

	RouteStoreErrorCode int
)

func (rsl routeStoreLocal) GetByCurrDest(curr Id, dest Id) (Route, StoreError) {
	routes, err := rsl.All()
	if err != nil {
		return Route{}, err
	}
	panic(routes)

}

func (rsl routeStoreLocal) All() (map[Id]Route, StoreError) {
	return rsl, nil
}
func (rsl routeStoreLocal) GetById(id Id) (Route, StoreError) {
	route, found := rsl[id]

	if !found {
		return Route{}, IdDoesntExist{}
	}

	return route, nil
}

func (rsl routeStoreLocal) Delete(id Id) StoreError {
	_, err := rsl.GetById(id)
	if err != nil {
		return err
	}

	// TODO: Actually delete the route

	return nil

}
