package main

import "github.com/google/uuid"

// Coordinate train and station schedules

type (
	trip struct {
		id                         id
		fromStationId, toStationId id
		trainId                    id
	}

	tripCoordinator interface {
		scheduleTrip(t trip) error
		cancelTrip(tripId id) error
	}

	tripStoreLocal struct {
		trips map[id]trip
	}
)

func newTrip(from, to id, train train) trip {
	// NOTE: not sure if this function would be better taking a station or just the id
	return trip{
		id:            uuid.New(),
		fromStationId: from,
		toStationId:   to,
		trainId:       train.E.Id,
	}
}

func newTripStoreLocal() *tripStoreLocal {
	return &tripStoreLocal{trips: map[id]trip{}}
}
