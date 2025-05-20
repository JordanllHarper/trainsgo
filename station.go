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

func (ssl stationStoreLocal) GetById(id Id) (Station, StoreError) {
	item, found := ssl.stations[id]
	if !found {
		return Station{}, IdDoesntExist(id)

	}

	return item, nil
}

func (ssl stationStoreLocal) All() (map[Id]Station, StoreError) {
	return maps.Clone(ssl.stations), nil
}

func (ssl *stationStoreLocal) getByName(name string) ([]Station, StoreError) {
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

type (
	registerStationErrorCode int

	stationPositionTaken Id
)

const (
	registerStationErrIdExists      registerStationErrorCode = 0
	registerStationErrPositionTaken registerStationErrorCode = 1
)

func (e idAlreadyExists) RegisterCode() registerStationErrorCode { return registerStationErrIdExists }

func (e stationPositionTaken) HttpCode() int { return http.StatusBadRequest }

func (e stationPositionTaken) RegisterCode() registerStationErrorCode {
	return registerStationErrPositionTaken
}

func (e stationPositionTaken) Error() string {
	return fmt.Sprintf("Station Position already taken by %s", Id(e))
}

func (ssl *stationStoreLocal) register(s Station) HttpError {
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

func (ssl stationStoreLocal) Delete(id Id) StoreError {
	// TODO: Cancel all schedules going to this station
	return nil
}

func (ssl *stationStoreLocal) changeName(id Id, newName string) StoreError {
	station, found := ssl.stations[id]
	if !found {
		return IdDoesntExist(id)
	}

	station.Name = newName
	ssl.stations[id] = station

	return nil
}
