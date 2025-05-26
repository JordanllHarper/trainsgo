package main

import (
	"cmp"
	"context"
	"fmt"
	"maps"
	"net/http"
	"slices"

	"github.com/google/uuid"
)

type (
	/*
		Represents a path (a route) from a CurrentId to a DestId.

		NextId represents the next line to take to get to DestId
	*/
	Route struct {
		Id        Id
		CurrentId Id
		DestId    Id
		NextId    Id
	}

	Set[T comparable] map[T]bool
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
		StoreDeleter[Route]
		GetByCurrDest(current, dest Id) (Id, error)
		WriteRoute(r Route) error
		WriteBatch(r []Route) error
		MapRoute(currentId, destId Id) ([]Line, error)
	}

	routeStoreLocal struct {
		routes   map[Id]Route
		stations IDable[Station]
	}
)

func NewRouteStoreLocal(stations StoreIDable[Station]) routeStoreLocal {
	return routeStoreLocal{
		map[Id]Route{},
		stations,
	}
}

func (rsl routeStoreLocal) GetByCurrDest(curr Id, dest Id) (Id, error) {
	wrapErr := func(err error) error {
		return fmt.Errorf("Get by Curr Dest %w", err)
	}
	routes, stErr := rsl.All()
	if stErr != nil {
		return Id{}, internalError{stErr}
	}

	for _, r := range routes {
		if r.CurrentId == curr && r.DestId == dest {
			return r.NextId, nil
		}
	}

	lineRoute, rtErr := rsl.MapRoute(curr, dest)
	if rtErr != nil {
		return Id{}, rtErr
	}

	if len(lineRoute) == 0 {
		return Id{}, wrapErr(fmt.Errorf("No routes"))
	}

	next := lineRoute[0].Id

	return next, nil
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

func (e idDoesntExist) Id() Id { return Id(e) }

type (
	routerChans struct {
		addCh   chan<- chan idSet
		visitCh chan<- Id
		errCh   chan<- error
		endCh   chan<- []Line
	}

	idSet Set[Id]

	noRouteFound struct {
		curr, dest Id
	}
	routerIdInvalidConnection Id
)

type searcher struct {
	id    Id
	updCh chan idSet
}

func (ri routeStoreLocal) MapRoute(currentStationId, destStationId Id) ([]Line, error) {
	visited := idSet{}
	searchers := map[Id]chan idSet{}
	routes := [][]Line{}

	updateVisited := func(
		ctx context.Context,
		visitUpdate chan Id,
	) {
		for {
			select {
			case visit := <-visitUpdate:
				visited[visit] = true
				for _, r := range searchers {
					r <- visited
				}

			case <-ctx.Done():
				return
			}
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	routeCh := make(chan []Line)

	visitCh := make(chan Id)
	go updateVisited(ctx, visitCh)

	searcherCh := make(chan searcher)
	rmvCh := make(chan Id)

	bfs(
		ri.stations,
		currentStationId,
		destStationId,
		[]Line{},
		routeCh,
		visitCh,
		searcherCh,
		rmvCh,
	)
loop:
	for {
		select {
		case s := <-searcherCh:
			fmt.Println("New Searcher", s.id)
			searchers[s.id] = s.updCh

		case r := <-routeCh:
			routes = append(routes, r)
		case id := <-rmvCh:
			fmt.Println("Removing Searcher", id)
			delete(searchers, id)
			if len(searchers) == 0 {
				break loop
			}
		}
	}
	smallestRoute := slices.MinFunc(routes, func(r1, r2 []Line) int {
		return cmp.Compare(len(r1), len(r2))
	})

	return smallestRoute, nil
}

func bfs(
	stations IDable[Station],
	currentStationId, destStationId Id,
	currentRoute []Line,
	endCh chan<- []Line,
	visitCh chan<- Id,
	searchCh chan<- searcher,
	rmvCh chan<- Id,
) {
	{
		if currentStationId == destStationId {
			endCh <- currentRoute
			return
		}
	}
	// Check if this is a valid station
	curr, err := stations.GetById(currentStationId)
	if err != nil {
		// TODO: Handle properly
		panic(err)
	}

	// setup the visited nodes for this searcher as well as the update function
	visited := idSet{}
	awaitIdSetUpdates := func(
		ctx context.Context,
		updVisited chan idSet,
	) {
		for {
			select {
			case s := <-updVisited:
				visited = s
			case <-ctx.Done():
				return
			}
		}
	}

	// Receive updates from the main goroutine to update visited
	updateRecCh := make(chan idSet)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go awaitIdSetUpdates(ctx, updateRecCh)

	// Register a new searcher
	searcherId := uuid.New()
	searchCh <- searcher{
		searcherId,
		updateRecCh,
	}
	// we will remove it as soon as we've spun up all the other routines
	defer func() { rmvCh <- searcherId }()

	// Update all goroutines with this station id
	visitCh <- currentStationId

	for _, line := range curr.SurroundingLines {
		if visited[line.Id] {
			continue
		}

		nextStId, err := getNonCurrent(currentStationId, line)
		if err != nil {
			// TODO: Handle
			panic(err)
		}

		go bfs(
			stations,
			nextStId,
			destStationId,
			append(currentRoute, line),
			endCh,
			visitCh,
			searchCh,
			rmvCh,
		)
	}
}

func (e routerIdInvalidConnection) Error() string {
	return fmt.Sprintf("ID %s invalid connecting station", Id(e))
}
func (e routerIdInvalidConnection) HttpCode() int { return http.StatusBadRequest }

func (e noRouteFound) Error() string {
	return fmt.Sprintf("No route found between %v - %v", e.curr, e.dest)
}
func (e noRouteFound) HttpCode() int { return http.StatusBadRequest }

func getNonCurrent(currentId Id, line Line) (Id, error) {
	var nextSt Id
	st1Id, st2Id := line.StationOne, line.StationTwo
	switch currentId {
	case st1Id:
		nextSt = st2Id
	case st2Id:
		nextSt = st1Id
	default:
		return Id{}, routerIdInvalidConnection(currentId)
	}
	return nextSt, nil
}
