package main

import "github.com/google/uuid"

type (
	tripStatus string
	Trip       struct {
		Id            id            `json:"id"`
		FromStationId id            `json:"fromStationId"`
		ToStationId   id            `json:"toStationId"`
		TrainId       id            `json:"trainId"`
		ExpectedTimes expectedTimes `json:"expectedTimes"`
		Status        tripStatus    `json:"status"`
	}

	tripCoordinator interface {
		scheduleTrip(t Trip) error
		delayTrip(id id) error
		cancelTrip(tripId id) error
	}

	tripReaderCoordinator interface {
		storeReader[Trip]
		tripCoordinator
	}

	tripCoordinatorLocal struct {
		trips    map[id]Trip
		trains   storeReader[Train]
		stations storeReader[Station]
	}
)

const (
	OnTime    tripStatus = "On time"
	Delayed              = "Delayed"
	Cancelled            = "Cancelled"
)

func newTrip(from, to id, train Train) Trip {
	// NOTE: not sure if this function would be better taking a station or just the id
	return Trip{
		Id:            uuid.New(),
		FromStationId: from,
		ToStationId:   to,
		TrainId:       train.E.Id,
	}
}

func newTripCoordinatorLocal(trains storeReader[Train], stations storeReader[Station]) *tripCoordinatorLocal {
	return &tripCoordinatorLocal{
		trips:    map[id]Trip{},
		trains:   trains,
		stations: stations,
	}
}

func (tcl *tripCoordinatorLocal) all() (map[id]Trip, error) {
	return tcl.trips, nil
}

func (tcl *tripCoordinatorLocal) getById(id id) (Trip, error) {
	t, found := tcl.trips[id]

	if !found {
		return Trip{}, newStoreReaderError(id, "Trip", StoreReaderErrIdNotFound)
	}

	return t, nil
}

func (tcl *tripCoordinatorLocal) delayTrip(id id) error {
	t, found := tcl.trips[id]

	if !found {
		return newStoreReaderError(id, "Trip", StoreReaderErrIdNotFound)
	}

	t.Status = Delayed

	tcl.trips[id] = t

	return nil
}

func (tcl *tripCoordinatorLocal) scheduleTrip(t Trip) error {
	// _, found := tcl.trips[t.Id]
	//
	// if found {
	// 	return newStoreReaderError(t.Id, "Trip", StoreReaderErrIdNotFound)
	// }
	//
	// train, err := tcl.trains.getById(t.TrainId)
	// if err != nil {
	// 	return err
	// }
	//
	// st1, err := tcl.stations.getById(t.FromStationId)
	// if err != nil {
	// 	return err
	// }
	//
	// st2, err := tcl.stations.getById(t.ToStationId)
	// if err != nil {
	// 	return err
	// }
	//
	// // Add the trip to the store
	// tcl.trips[t.Id] = t
	//
	// return nil
	return nil
}

func (tcl *tripCoordinatorLocal) cancelTrip(tripId id) error {
	return nil
}
