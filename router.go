package main

import (
	"fmt"
	"net/http"
)

const (
	RouterErrIdDoesntExist     RouterErrorCode = 0
	RouterErrInvalidConnection RouterErrorCode = 1
	RouterErrInternalError     RouterErrorCode = 2
)

type (
	/*
		Given a currentId for the current location and a destinationId, what is the next line the user needs to take.
	*/
	Router interface {
		Route(currentId Id, destId Id) (Route, RouterError)
	}

	routerImpl struct {
		stations StoreReader[Station]
		routes   StoreReader[Route]
	}

	routerChans struct {
		addCh     chan chan idSet
		visitedCh chan idSet
		errCh     chan RouterError
		endCh     chan []Line
	}

	idSet map[Id]bool

	RouterErrorCode int
	RouterError     interface {
		HttpError
		RouterErrorCode() RouterErrorCode
	}
	routerIdInvalidConnection Id
)

func (e idDoesntExist) RouterErrorCode() RouterErrorCode       { return RouterErrIdDoesntExist }
func (e internalServerError) RouterErrorCode() RouterErrorCode { return RouterErrInternalError }
func (e routerIdInvalidConnection) Error() string {
	return fmt.Sprintf("ID %s invalid connecting station", Id(e))
}
func (e routerIdInvalidConnection) HttpCode() int { return http.StatusBadRequest }
func (e routerIdInvalidConnection) RouterErrorCode() RouterErrorCode {
	return RouterErrInvalidConnection
}

func (ri routerImpl) Route(currentId Id, destId Id) (Route, RouterError) {

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
	errCh := make(chan RouterError)
	// channel for subscribing a new routine to updates
	addCh := make(chan chan idSet)
	// all the channels that should receive updates from about the visited nodes
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
	}

	// update the channels as we get requests
	for {
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
		}
	}

}

func getStation(stations StoreReader[Station], currentId Id) (Station, RouterError) {
	st, stErr := stations.GetById(currentId)
	if stErr != nil {
		switch stErr.StoreErrorCode() {
		case StoreErrorIdDoesntExist:
			return Station{}, idDoesntExist(currentId)
		case StoreErrorInternalError:
			return Station{}, internalServerError{stErr}
		}
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
	errChan chan RouterError,
	endChan chan []Line,
) routerChans {
	return routerChans{addChan, visitedChan, errChan, endChan}
}
