package main

import (
	"fmt"

	"github.com/google/uuid"
)

type (
	TripStatus int
	Trip       struct {
		Id            Id            `json:"id"`
		FromStationId Id            `json:"fromStationId"`
		ToStationId   Id            `json:"toStationId"`
		TrainId       Id            `json:"trainId"`
		ExpectedTimes ExpectedTimes `json:"expectedTimes"`
		Status        TripStatus    `json:"status"`
	}

	tripHandlerLocal struct {
		trips    tripStoreLocal
		trains   StoreReader[Train]
		stations StoreReader[Station]
		router   Router
	}
	tripStoreLocal map[Id]Trip
)

const (
	OnTime    TripStatus = 0
	Delayed              = 1
	Cancelled            = 2
)

func (ts TripStatus) String() string {
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

func NewTrip(from, to Station, train Train, expTimes ExpectedTimes, status TripStatus) Trip {
	return Trip{
		uuid.New(),
		from.E.Id,
		to.E.Id,
		train.E.Id,
		expTimes,
		status,
	}
}

func NewTripCoordinatorLocal(trains StoreReader[Train], stations StoreReader[Station]) *tripHandlerLocal {
	return &tripHandlerLocal{
		trips:    map[Id]Trip{},
		trains:   trains,
		stations: stations,
	}
}

func (tcl tripStoreLocal) All() (map[Id]Trip, StoreError) {
	return tcl, nil
}

func (tsl tripStoreLocal) GetById(id Id) (Trip, StoreError) {
	t, found := tsl[id]

	if !found {
		return Trip{}, IdDoesntExist(id)
	}

	return t, nil
}

func (tcl *tripHandlerLocal) delayTrip(id Id) error {
	t, found := tcl.trips[id]

	if !found {
		return idDoesntExist(id)
	}

	t.Status = Delayed

	tcl.trips[id] = t

	return nil
}

func (tcl *tripHandlerLocal) cancelTrip(tripId Id) error {
	return nil
}
