package main

import (
	"cmp"
	"fmt"
	"net/http"
	"slices"
)

type (
	/*
		Given a currentId for the current location and a destinationId, what is the next line the user needs to take.
	*/

	routerChans struct {
		addCh     chan chan idSet
		visitedCh chan idSet
		errCh     chan error
		endCh     chan []Line
	}

	idSet map[Id]bool

	noRouteFound struct {
		curr, dest Id
	}
	routerIdInvalidConnection Id
)

func (ri routeStoreLocal) MapRoute(currentId Id, destId Id) (Route, error) {
	// The found routes
	routes := [][]Line{}
	// The ids that have been visited
	visitedIds := idSet{}
	// chans
	// When a node is visited, this will receive a value
	visitCh := make(chan Id)
	// For when a route has been found to the destId
	endCh := make(chan []Line)
	// When a channel reports an error
	errCh := make(chan error)
	// channel for subscribing a new routine to updates
	addCh := make(chan chan idSet)
	// all the channels that should receive updates about the visited nodes
	routineChs := []chan idSet{}

	root, err := getStation(ri.stations, currentId)
	if err != nil {
		return Route{}, err
	}

	surroundingLines := root.SurroundingLines
	for _, line := range surroundingLines {
		routineVisitedCh := make(chan idSet)
		go bfs(
			ri.stations,
			[]Line{line},
			line.Id,
			destId,
			newRouterChans(
				addCh,
				routineVisitedCh,
				errCh,
				endCh,
			),
		)
		routineChs = append(routineChs, routineVisitedCh)
	}

	// update the channels as we get requests
	numRoutines := len(routineChs)
	check := true
	for check {
		select {
		case err := <-errCh:
			// TODO: We could have some resilience here and continue to search, just cancel the routine that ended
			panic(err)
		case id := <-visitCh:
			visitedIds[id] = true
			for _, ch := range routineChs {
				ch <- visitedIds
			}
		case r := <-endCh:
			routes = append(routes, r)
		case ch := <-addCh:
			routineChs = append(routineChs, ch)
			numRoutines++
		default:
			if numRoutines == 0 {
				// We have finished
				check = false
			}
		}
	}
	if len(routes) == 0 {
		return Route{}, noRouteFound{currentId, destId}
	}
	minLine := slices.MinFunc(routes, func(left, right []Line) int {
		return cmp.Compare(len(left), len(right))
	})
	if len(minLine) == 0 {
		return Route{}, noRouteFound{currentId, destId}
	}

	route := NewRoute(currentId, destId, minLine[0].Id)
	err = ri.WriteRoute(route)
	if err != nil {
		return Route{}, err
	}
	ids, err := writeRoutes(ri, minLine[1:], currentId, destId)
	if err != nil {
		if e := ri.DeleteBatch(ids); e != nil {
			return Route{}, internalError{e}
		}

		return Route{}, err
	}
	return route, nil
}

func writeRoutes(r RouteStore, lines []Line, currentId, destId Id) ([]Id, error) {
	routeIdsWritten := []Id{}
	for _, l := range lines {
		route := NewRoute(currentId, destId, l.Id)
		err := r.WriteRoute(route)
		if err != nil {
			return routeIdsWritten, err
		}
		routeIdsWritten = append(routeIdsWritten, route.Id)
	}
	return routeIdsWritten, nil
}

func getStation(stations StoreReader[Station], currentId Id) (Station, error) {
	st, err := stations.GetById(currentId)
	if err != nil {
		return Station{}, err
	}
	return st, nil
}

func bfs(
	stations StoreReader[Station],
	route []Line,
	currentId Id,
	destId Id,
	chans routerChans,
) {
	// if we have reached our id, we're done
	if currentId == destId {
		chans.endCh <- route
		return
	}

	st, err := getStation(stations, currentId)
	if err != nil {
		// cancel this search
		// TODO: We could entertain continuing here
		chans.errCh <- err
		return
	}

	// our main line iteration loop
	go func() {
		visited := map[Id]bool{}
		visitedUpdates := make(chan idSet)
		chans.addCh <- visitedUpdates
		go func() {
			for {
				select {
				case update, more := <-visitedUpdates:
					if more {
						visited = update
						continue
					}
					// cancel this channel if the visitedUpdates channel is closed
					close(visitedUpdates)
				}
			}
		}()
		for _, line := range st.SurroundingLines {
			st1Id, st2Id := line.StationOne, line.StationOne
			var nextSt Id
			switch currentId {
			case st1Id:
				nextSt = st2Id
			case st2Id:
				nextSt = st1Id
			default:
				chans.errCh <- routerIdInvalidConnection(currentId)
				return
			}
			if visited[nextSt] {
				continue
			}
			{
				visited[nextSt] = true
				chans.visitedCh <- visited
			}
			st, err := getStation(stations, nextSt)
			if err != nil {
				chans.errCh <- idDoesntExist(st.E.Id)
				return
			}
			go bfs(stations, route, st.E.Id, destId, chans)
		}
	}( // We get the station connected to this one, which is NOT this one (since we have no guarantees of which station is which)
	)

}

func newRouterChans(
	addChan chan chan idSet,
	visitedChan chan idSet,
	errChan chan error,
	endChan chan []Line,
) routerChans {
	return routerChans{addChan, visitedChan, errChan, endChan}
}

func (e routerIdInvalidConnection) Error() string {
	return fmt.Sprintf("ID %s invalid connecting station", Id(e))
}
func (e routerIdInvalidConnection) HttpCode() int { return http.StatusBadRequest }

func (e noRouteFound) Error() string {
	return fmt.Sprintf("No route found between %v - %v", e.curr, e.dest)
}
func (e noRouteFound) HttpCode() int { return http.StatusBadRequest }
