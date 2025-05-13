package main

import (
	"fmt"
	"slices"
)

type (
	registrar[T any] interface {
		register(s T) error
		deregister(id id) error
	}

	store[T any] interface {
		all() ([]T, error)
		getById(id id) (T, error)
	}
)

// Local store implementation

type (
	stationStoreLocal struct {
		stations []station
	}

	trainStoreLocal struct {
		trains []train
	}

	lineStoreLocal struct {
		lines []line
	}
)

func (ssl *stationStoreLocal) getById(id id) (station, error) {
	item, found := sliceGet(
		ssl.stations,
		func(s station) bool {
			return s.id == id
		},
	)

	if !found {
		return station{},
			errorIdNotFound(id, "Station")
	}

	return item, nil
}

func (ssl *stationStoreLocal) all() ([]station, error) {
	return ssl.stations, nil
}

func (ssl *stationStoreLocal) register(s station) error {
	if slices.ContainsFunc(
		ssl.stations,
		func(s2 station) bool {
			return s.id == s2.id
		}) {
		return idAlreadyExists(s.id, "Station")
	}

	if slices.ContainsFunc(
		ssl.stations,
		func(s2 station) bool {
			return s.position == s2.position
		}) {
		return fmt.Errorf(
			"There is already a Station at position %s",
			s.position,
		)
	}

	ssl.stations = append(ssl.stations, s)

	return nil
}

func (ssl *stationStoreLocal) deregister(id id) error {
	// TODO: Cancel all schedules going to this station
	return nil
}

//

func (tsl *trainStoreLocal) all() ([]train, error) {
	return tsl.trains, nil
}

func (tsl *trainStoreLocal) getById(id id) (train, error) {
	t, found := sliceGet(
		tsl.trains,
		func(t train) bool {
			return t.id == id
		},
	)

	if !found {
		return train{},
			errorIdNotFound(id, "Train")
	}

	return t, nil
}

func (tsl *trainStoreLocal) register(t train) error {
	if slices.ContainsFunc(tsl.trains, func(t2 train) bool {
		return t.id == t2.id
	}) {
		return idAlreadyExists(t.id, "Train")
	}

	tsl.trains = append(tsl.trains, t)

	return nil
}

func (tsl trainStoreLocal) deregister(id id) error {
	// TODO: Finish this trains schedule and then remove
	return nil
}

//

func (lsl *lineStoreLocal) all() ([]line, error) {
	return lsl.lines, nil
}

func (lsl *lineStoreLocal) getById(id id) (line, error) {
	value, found := sliceGet(lsl.lines, func(l line) bool {
		return l.id == id
	})

	if !found {
		return line{}, errorIdNotFound(id, "Line")
	}
	return value, nil
}

func (lsl *lineStoreLocal) register(l line) error {
	lsl.lines = append(lsl.lines, l)
	return nil
}

func (lsl *lineStoreLocal) deregister(id id) error {
	// TODO: Wait for trains to finish using this line, then decommission
	return nil
}
