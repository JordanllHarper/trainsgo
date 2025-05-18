package main

import (
	"fmt"
	"maps"
)

type (
	Station struct {
		E                entity `json:"entity"`
		Name             string `json:"name"`
		Platforms        int    `json:"platforms"`
		SurroundingLines []Line `json:"surroundingLines"`
	}

	errStationAlreadyAtPosition struct {
		id
		position
	}

	stationStoreLocal struct {
		stations map[id]Station
	}
)

func newStationStoreLocal() *stationStoreLocal {
	return &stationStoreLocal{stations: map[id]Station{}}
}

func newStation(pos position, name string, platforms int) Station {
	return Station{
		E:         newEntity(pos),
		Name:      name,
		Platforms: platforms,
	}
}

func (s Station) String() string {
	return fmt.Sprintf(
		"{ID: %v, Name: %v, Platforms: %v}",
		s.E.Id,
		s.Name,
		s.Platforms,
	)
}

func (ssl stationStoreLocal) getById(id id) (Station, *storeReaderError) {
	item, found := ssl.stations[id]
	if !found {
		return Station{},
			newStoreReaderError(id, "Station", StoreReaderErrIdNotFound)
	}

	return item, nil
}

func (ssl stationStoreLocal) all() (map[id]Station, *storeReaderError) {
	return maps.Clone(ssl.stations), nil
}

func (ssl *stationStoreLocal) getByName(name string) ([]Station, *storeReaderError) {
	stations := []Station{}
	for v := range maps.Values(ssl.stations) {
		if v.Name == name {
			stations = append(stations, v)
		}
	}

	return stations, nil
}

func (err errStationAlreadyAtPosition) Error() string {
	return fmt.Sprintf(
		"There is already Station %v at position %s",
		err.id,
		err.position,
	)
}

type registerStationErrorCode int

const (
	registerStationErrIdExists registerStationErrorCode = iota
	registerStationErrPositionTaken
)

type registerStationError struct {
	id   id
	code registerStationErrorCode
}

func (ssl *stationStoreLocal) register(s Station) *registerStationError {
	_, found := ssl.stations[s.E.Id]

	if found {
		return &registerStationError{s.E.Id, registerStationErrIdExists}
	}

	for v := range maps.Values(ssl.stations) {
		if v.E.Pos == s.E.Pos {
			return &registerStationError{s.E.Id, registerStationErrPositionTaken}
		}
	}

	ssl.stations[s.E.Id] = s

	return nil
}

func (ssl stationStoreLocal) delete(id id) *storeDeleterError {
	// TODO: Cancel all schedules going to this station
	return nil
}

func (ssl *stationStoreLocal) changeName(id id, newName string) *storeReaderError {
	station, found := ssl.stations[id]
	if !found {
		return newStoreReaderError(id, "Station", StoreReaderErrIdNotFound)
	}

	station.Name = newName
	ssl.stations[id] = station

	return nil
}
