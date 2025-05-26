package main

import (
	"fmt"

	"github.com/google/uuid"
)

const (
	OnTime    TripStatus = 0
	Delayed   TripStatus = 1
	Cancelled TripStatus = 2
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
		trains   StoreIDable[Train]
		stations StoreIDable[Station]
		router   RouteStore
	}
	tripStoreLocal map[Id]Trip
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

func (tcl tripStoreLocal) All() (map[Id]Trip, error) {
	return tcl, nil
}

func (tsl tripStoreLocal) GetById(id Id) (Trip, error) {
	t, found := tsl[id]

	if !found {
		return Trip{}, idDoesntExist(id)
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
