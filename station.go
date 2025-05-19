package main

import (
	"fmt"
	"maps"
)

type (
	Station struct {
		E                Entity `json:"entity"`
		Name             string `json:"name"`
		Platforms        int    `json:"platforms"`
		SurroundingLines []Line `json:"surroundingLines"`
	}

	errStationAlreadyAtPosition struct {
		Id
		Position
	}

	stationStoreLocal struct {
		stations map[Id]Station
	}
)

func NewStationStoreLocal() *stationStoreLocal {
	return &stationStoreLocal{stations: map[Id]Station{}}
}

func NewStation(pos Position, name string, platforms int) Station {
	return Station{
		E:         NewEntity(pos),
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

func (ssl stationStoreLocal) GetById(id Id) (Station, *StoreReaderError) {
	item, found := ssl.stations[id]
	if !found {
		return Station{},
			NewStoreReaderError(id, "Station", StoreReaderErrIdNotFound)
	}

	return item, nil
}

func (ssl stationStoreLocal) All() (map[Id]Station, *StoreReaderError) {
	return maps.Clone(ssl.stations), nil
}

func (ssl *stationStoreLocal) getByName(name string) ([]Station, *StoreReaderError) {
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
		err.Id,
		err.Position,
	)
}

type registerStationErrorCode int

const (
	registerStationErrIdExists registerStationErrorCode = iota
	registerStationErrPositionTaken
)

type registerStationError struct {
	id   Id
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

func (ssl stationStoreLocal) Delete(id Id) *StoreDeleterError {
	// TODO: Cancel all schedules going to this station
	return nil
}

func (ssl *stationStoreLocal) changeName(id Id, newName string) *StoreReaderError {
	station, found := ssl.stations[id]
	if !found {
		return NewStoreReaderError(id, "Station", StoreReaderErrIdNotFound)
	}

	station.Name = newName
	ssl.stations[id] = station

	return nil
}
