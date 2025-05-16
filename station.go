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

	stationStoreLocal struct {
		stations map[id]Station
	}
)

func newStationStoreLocal() *stationStoreLocal {
	return &stationStoreLocal{stations: map[id]Station{}}
}

func (ssl *stationStoreLocal) changeName(id id, newName string) error {
	station, found := ssl.stations[id]
	if !found {
		return newErrIdNotFound(id, "Station")
	}

	station.Name = newName
	ssl.stations[id] = station

	return nil
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

func (ssl *stationStoreLocal) getById(id id) (Station, error) {
	item, found := ssl.stations[id]
	if !found {
		return Station{},
			newErrIdNotFound(id, "Station")
	}

	return item, nil
}

func (ssl *stationStoreLocal) all() (map[id]Station, error) {
	return maps.Clone(ssl.stations), nil
}

func (ssl *stationStoreLocal) getByName(name string) ([]Station, error) {
	stations := []Station{}
	for v := range maps.Values(ssl.stations) {
		if v.Name == name {
			stations = append(stations, v)
		}
	}

	return stations, nil
}

type errStationAlreadyAtPosition struct {
	id
	position
}

func (err errStationAlreadyAtPosition) Error() string {
	return fmt.Sprintf(
		"There is already Station %v at position %s",
		err.id,
		err.position,
	)
}

func newErrStationAlreadyAtPosition(id id, pos position) errStationAlreadyAtPosition {
	return errStationAlreadyAtPosition{id, pos}
}

func (ssl *stationStoreLocal) register(s Station) error {

	_, found := ssl.stations[s.E.Id]

	if found {
		return newErrIdAlreadyExists(s.E.Id, "Station")
	}

	for v := range maps.Values(ssl.stations) {
		if v.E.Pos == s.E.Pos {
			return newErrStationAlreadyAtPosition(s.E.Id, s.E.Pos)
		}
	}

	ssl.stations[s.E.Id] = s

	return nil
}

func (ssl *stationStoreLocal) deregister(id id) error {
	// TODO: Cancel all schedules going to this station
	return nil
}
