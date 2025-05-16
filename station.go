package main

import (
	"fmt"
	"maps"
)

type (
	station struct {
		entity           `json:"entity"`
		Name             string `json:"name"`
		Platforms        int    `json:"platforms"`
		SurroundingLines []line `json:"surrounding_lines"`
	}

	stationStoreLocal struct {
		stations map[id]station
	}
)

func newStationStoreLocal() *stationStoreLocal {
	return &stationStoreLocal{stations: map[id]station{}}
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

func newStation(pos position, name string, platforms int) station {
	return station{
		entity:    newEntity(pos),
		Name:      name,
		Platforms: platforms,
	}
}

func (s station) String() string {
	return fmt.Sprintf(
		"{ID: %v, Name: %v, Platforms: %v}",
		s.id,
		s.Name,
		s.Platforms,
	)
}

func (ssl *stationStoreLocal) getById(id id) (station, error) {
	item, found := ssl.stations[id]
	if !found {
		return station{},
			newErrIdNotFound(id, "Station")
	}

	return item, nil
}

func (ssl *stationStoreLocal) all() (map[id]station, error) {
	return maps.Clone(ssl.stations), nil
}

func (ssl *stationStoreLocal) getByName(name string) ([]station, error) {
	stations := []station{}
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

func (ssl *stationStoreLocal) register(s station) error {

	_, found := ssl.stations[s.id]

	if found {
		return newErrIdAlreadyExists(s.id, "Station")
	}

	for v := range maps.Values(ssl.stations) {
		if v.position == s.position {
			return newErrStationAlreadyAtPosition(s.id, s.position)
		}
	}

	ssl.stations[s.id] = s

	return nil
}

func (ssl *stationStoreLocal) deregister(id id) error {
	// TODO: Cancel all schedules going to this station
	return nil
}
