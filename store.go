package main

import (
	"errors"
	"fmt"
	"slices"
)

type (
	stationStore interface {
		getAll() ([]station, error)
		getById(id id) (station, error)
		register(s station) error
		deregister(id id) error
	}

	trainStore interface {
		getAll() ([]train, error)
		getById(id id) (train, error)
		register(t train) error
		deregister(id id) error
	}

	lineStore interface {
		getAll() ([]line, error)
		register(l line) error
		deregister(id id) error
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
		return station{}, errors.New(
			fmt.Sprintf(
				"Station with ID %s doesn't exist",
				id,
			),
		)
	}

	return item, nil
}

func (ssl *stationStoreLocal) getAll() ([]station, error) {
	return ssl.stations, nil
}

func (ssl *stationStoreLocal) register(s station) error {
	if slices.ContainsFunc(
		ssl.stations,
		func(s2 station) bool {
			return s.id == s2.id
		}) {
		return errors.New(
			fmt.Sprintf(
				"Registered station id %s already exists",
				s.id,
			),
		)
	}

	if slices.ContainsFunc(
		ssl.stations,
		func(s2 station) bool {
			return s.position == s2.position
		}) {
		return errors.New(
			fmt.Sprintf(
				"There is already a station at position %s",
				s.position,
			),
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

func (tsl *trainStoreLocal) getAll() ([]train, error) {
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
		return train{}, errors.New(
			fmt.Sprintf(
				"Train with ID %s doesn't exist",
				id,
			),
		)
	}

	return t, nil
}

func (tsl *trainStoreLocal) register(t train) error {
	if slices.ContainsFunc(tsl.trains, func(t2 train) bool {
		return t.id == t2.id
	}) {

		return errors.New(
			fmt.Sprintf(
				"Registered train ID %s already exists",
				t.id,
			),
		)
	}

	tsl.trains = append(tsl.trains, t)

	return nil
}

func (tsl trainStoreLocal) deregister(id id) error {
	// TODO: Finish this trains schedule and then remove
	return nil
}

//

func (lsl *lineStoreLocal) getAll() ([]line, error) {
	return lsl.lines, nil
}

func (lsl *lineStoreLocal) register(l line) error {
	lsl.lines = append(lsl.lines, l)
	return nil
}

func (lsl *lineStoreLocal) deregister(id id) error {
	// TODO: Wait for trains to finish using this line, then decommission
	return nil
}
