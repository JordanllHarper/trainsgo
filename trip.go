package main

import (
	"fmt"

	"github.com/google/uuid"
)

type (
	tripStatus int
	Trip       struct {
		Id            Id            `json:"id"`
		FromStationId Id            `json:"fromStationId"`
		ToStationId   Id            `json:"toStationId"`
		TrainId       Id            `json:"trainId"`
		ExpectedTimes expectedTimes `json:"expectedTimes"`
		Status        tripStatus    `json:"status"`
	}

	tripHandlerLocal struct {
		trips    tripStoreLocal
		trains   storeReader[Train]
		stations storeReader[Station]
	}
	tripStoreLocal map[Id]Trip
)

const (
	OnTime    tripStatus = 0
	Delayed              = 1
	Cancelled            = 2
)

func (ts tripStatus) String() string {
	switch ts {
	case OnTime:
		return "On Time"
	case Cancelled:
		return "Cancelled"
	case Delayed:
		return "Delayed"
	default:
		panic(fmt.Sprintf("unexpected main.tripStatus: %#v", ts))
	}
}

func newTrip(from, to Station, train Train, expTimes expectedTimes, status tripStatus) Trip {
	return Trip{
		uuid.New(),
		from.E.Id,
		to.E.Id,
		train.E.Id,
		expTimes,
		status,
	}
}

func newTripCoordinatorLocal(trains storeReader[Train], stations storeReader[Station]) *tripHandlerLocal {
	return &tripHandlerLocal{
		trips:    map[Id]Trip{},
		trains:   trains,
		stations: stations,
	}
}

func (tcl tripStoreLocal) all() (map[Id]Trip, *storeReaderError) {
	return tcl, nil
}

func (tsl tripStoreLocal) getById(id Id) (Trip, *storeReaderError) {
	t, found := tsl[id]

	if !found {
		return Trip{}, newStoreReaderError(id, "Trip", StoreReaderErrIdNotFound)
	}

	return t, nil
}

func (tcl *tripHandlerLocal) delayTrip(id Id) error {
	t, found := tcl.trips[id]

	if !found {
		return newStoreReaderError(id, "Trip", StoreReaderErrIdNotFound)
	}

	t.Status = Delayed

	tcl.trips[id] = t

	return nil
}

func (tcl *tripHandlerLocal) scheduleTrip(t Trip) error {
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

func (tcl *tripHandlerLocal) cancelTrip(tripId Id) error {
	return nil
}
