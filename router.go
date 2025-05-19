package main

import (
	"fmt"
)

type (
	/*
		Given a currentId for the current location and a destinationId, what is the next set of lines the user needs to take.
	*/
	Router interface {
		Route(currentId Id, destId Id) ([]Line, *RouterError)
	}

	routerImpl struct {
		stations StoreReader[Station]
	}

	RouterErrorCode int
	RouterError     struct {
		Id   Id
		Code RouterErrorCode
	}

	RouterChans struct {
		addCh     chan chan idSet
		visitedCh chan idSet
		errCh     chan RouterError
		endCh     chan []Line
	}

	idSet map[Id]bool
)

func NewRouterChans(
	addChan chan chan idSet,
	visitedChan chan idSet,
	errChan chan RouterError,
	endChan chan []Line,
) RouterChans {
	return RouterChans{addChan, visitedChan, errChan, endChan}
}

const (
	RouterErrIdNotFound       RouterErrorCode = 0
	RouterErrIdNotAConnection RouterErrorCode = 1
	RouterErrInternalError    RouterErrorCode = 2
)

func (err RouterError) Error() string {
	switch err.Code {
	case RouterErrIdNotFound:
		return fmt.Sprintf("ID not found %s", err.Id)
	case RouterErrInternalError:
		return fmt.Sprintf("Internal Server Error")
	default:
		panic(fmt.Sprintf("unexpected main.RouterErrorCode: %#v", err.Code))
	}
}

func (ri routerImpl) Route(currentId Id, destId Id) ([]Line, *RouterError) {
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
		return nil, err
	}

	surroundingLines := root.SurroundingLines
	for _, line := range surroundingLines {
		routineVisitedCh := make(chan idSet)
		go bfs(
			ri.stations,
			[]Line{line},
			line.Id,
			destId,
			NewRouterChans(
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

func getStation(stations StoreReader[Station], currentId Id) (Station, *RouterError) {
	st, stErr := stations.GetById(currentId)
	if stErr != nil {
		switch stErr.code {
		case StoreReaderErrIdNotFound:
			return Station{}, &RouterError{Id: currentId, Code: RouterErrIdNotFound}
		case StoreReaderErrInternalError:
			return Station{}, &RouterError{Id: currentId, Code: RouterErrInternalError}
		default:
			panic(fmt.Sprintf("unexpected main.StoreReaderErrorCode: %#v", stErr.code))
		}
	}
	return st, nil
}

func bfs(
	stations StoreReader[Station],
	route []Line,
	currentId Id,
	destId Id,
	chans RouterChans,
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
		chans.errCh <- *err
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
			var connectingSt Id
			switch currentId {
			case st1Id:
				connectingSt = st2Id
			case st2Id:
				connectingSt = st1Id
			default:
				chans.errCh <- RouterError{Id: connectingSt, Code: RouterErrIdNotAConnection}
				return
			}
			if visited[connectingSt] {
				continue
			}
			{
				visited[connectingSt] = true
				chans.visitedCh <- visited
			}
			st, err := getStation(stations, connectingSt)
			if err != nil {
				chans.errCh <- RouterError{Id: connectingSt, Code: RouterErrIdNotFound}
				return
			}
			go bfs(stations, route, st.E.Id, destId, chans)
		}
	}( // We get the station connected to this one, which is NOT this one (since we have no guarantees of which station is which)
	)

}
