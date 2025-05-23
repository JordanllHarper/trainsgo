package main

import (
	"maps"
	"slices"

	"github.com/google/uuid"
)

type (
	/*
		Represents a path (a route) from a CurrentId to a DestId.

		NextId represents the next station or intersection or whatever to go to to continue the route.
	*/
	Route struct {
		Id        Id
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

type (
	RouteStore interface {
		StoreReaderDeleter[Route]
		GetByCurrDest(current, dest Id) (Route, error)
		WriteRoute(r Route) error
		WriteBatch(r []Route) error
		MapRoute(currentId Id, destId Id) (Route, error)
	}

	routeStoreLocal struct {
		routes   map[Id]Route
		stations StoreReader[Station]
	}
)

func (rsl routeStoreLocal) GetByCurrDest(curr Id, dest Id) (Route, error) {
	routes, stErr := rsl.All()
	if stErr != nil {
		return Route{}, internalError{stErr}
	}

	for _, r := range routes {
		if r.CurrentId == curr && r.DestId == dest {
			return r, nil
		}
	}

	route, rtErr := rsl.MapRoute(curr, dest)

	if rtErr != nil {
		return Route{}, rtErr
	}

	return route, nil
}

func (rsl routeStoreLocal) All() (map[Id]Route, error) { return rsl.routes, nil }

func (rsl routeStoreLocal) GetById(id Id) (Route, error) {
	route, found := rsl.routes[id]

	if !found {
		return Route{}, idDoesntExist(id)
	}

	return route, nil
}

func (rsl routeStoreLocal) WriteRoute(r Route) error {

	_, found := rsl.routes[r.Id]
	if found {
		return idAlreadyExists(r.Id)
	}

	rsl.routes[r.Id] = r

	return nil
}

func (rsl routeStoreLocal) WriteBatch(r []Route) error {
	ids := slices.Collect(maps.Keys(rsl.routes))
	for _, v := range r {
		if slices.Contains(ids, v.Id) {
			return idAlreadyExists(v.Id)
		}
		rsl.routes[v.Id] = v
	}
	return nil
}

func (rsl routeStoreLocal) Delete(id Id) error {
	_, err := rsl.GetById(id)
	if err != nil {
		return err
	}

	// TODO: Actually delete the route

	return nil

}
func (rsl routeStoreLocal) DeleteBatch(ids []Id) error {
	// TODO: Delete all routes
	return nil
}

func (e idDoesntExist) Id() Id {
	return Id(e)
}
