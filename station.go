package main

import (
	"fmt"
	"maps"
	"net/http"
)

type (
	Station struct {
		E                Entity `json:"entity"`
		Name             string `json:"name"`
		Platforms        int    `json:"platforms"`
		SurroundingLines []Line `json:"surroundingLines"`
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

func (ssl stationStoreLocal) GetById(id Id) (Station, error) {
	item, found := ssl.stations[id]
	if !found {
		return Station{}, idDoesntExist(id)

	}

	return item, nil
}

func (ssl stationStoreLocal) All() (map[Id]Station, error) {
	return ssl.stations, nil
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

type (
	stationPositionTaken Id
)

func (e stationPositionTaken) HttpCode() int { return http.StatusBadRequest }
func (e stationPositionTaken) Error() string {
	return fmt.Sprintf("Station Position already taken by %s", Id(e))
}

func (ssl stationStoreLocal) register(s Station) HttpError {
	_, found := ssl.stations[s.E.Id]

	if found {
		return idAlreadyExists(s.E.Id)
	}

	for v := range maps.Values(ssl.stations) {
		if v.E.Pos == s.E.Pos {
			return stationPositionTaken(v.E.Id)
		}
	}

	ssl.stations[s.E.Id] = s

	return nil
}

func (ssl stationStoreLocal) changeName(id Id, newName string) error {
	station, found := ssl.stations[id]
	if !found {
		return idDoesntExist(id)
	}

	station.Name = newName
	ssl.stations[id] = station

	return nil
}

func (ssl stationStoreLocal) Delete(id Id) error {
	// TODO: Cancel all schedules going to this station
	return nil
}

func (ssl stationStoreLocal) DeleteBatch(ids []Id) error {
	// TODO: Cancel all schedules going to stations
	return nil
}
