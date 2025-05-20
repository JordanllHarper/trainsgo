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

	routeStoreLocal struct {
		routes map[Id]Route
		router Router
	}

	RouteStoreErrorCode int
)

func (e IdDoesntExist) RouterErrorCode() RouterErrorCode {
	return RouterErrIdDoesntExist
}

func (rsl routeStoreLocal) GetByCurrDest(curr Id, dest Id) (Route, RouterError) {
	routes, stErr := rsl.All()
	if stErr != nil {
		switch stErr.StoreErrorCode() {
		case StoreErrorIdDoesntExist:
			return Route{}, idDoesntExist()
		case StoreErrorInternalError:
			return Route{}, internalServerError{stErr}
		}
	}

	for _, r := range routes {
		if r.CurrentId == curr && r.DestId == dest {
			return r, nil
		}
	}

	route, rtErr := rsl.router.Route(curr, dest)

	if rtErr != nil {
		return Route{}, rtErr
	}

	return route, nil
}

func (rsl routeStoreLocal) All() (map[Id]Route, StoreError) {
	return rsl.routes, nil
}
func (rsl routeStoreLocal) GetById(id Id) (Route, StoreError) {
	route, found := rsl.routes[id]

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
