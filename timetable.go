package main

import (
	"slices"
	"time"

	"github.com/google/uuid"
)

type (
	timetableViewer interface {
		all() ([]timetableEntry, error)
		getById(id id) (timetableEntry, error)
	}

	timetabler interface {
		add(e timetableEntry) error
		delay(by time.Duration) error
		cancelAndRemove(id id) error
	}

	timetableStatus int

	timetableEntry struct {
		id                 id
		trainId, stationId id
		platform           int
		timetableStatus
		expectedTimes
	}

	timetablerImpl struct {
		timetable []timetableEntry
	}
)

func (t timetablerImpl) all() ([]timetableEntry, error) {
	return t.timetable, nil
}

func (t timetablerImpl) getById(id id) (timetableEntry, error) {
	value, found := sliceGet(t.timetable, func(v timetableEntry) bool {
		return id == v.id
	})

	if !found {
		return timetableEntry{}, errorIdNotFound(id, "Timetable entry")
	}

	return value, nil
}

func (t timetablerImpl) add(e timetableEntry) error {
	if slices.ContainsFunc(t.timetable, func(e2 timetableEntry) bool {
		return e.id == e2.id
	}) {
		return idAlreadyExists(e.id, "Timetable Entry")
	}

	/*
		TODO: Prevent edge cases:
		- if the platform and time range is already in use, conflict!
		- think about this more ...
	*/

	t.timetable = append(t.timetable, e)

	return nil
}

func (t timetablerImpl) delay(by time.Duration) error {
	panic("not implemented") // TODO: Implement
}

func (t timetablerImpl) cancelAndRemove(id id) error {
	panic("not implemented") // TODO: Implement
}

const (
	Ok timetableStatus = iota
	Delayed
	Cancelled
)

func newTimetableEntry(trainId, stationId id, platform int, exp expectedTimes, status timetableStatus) timetableEntry {
	return timetableEntry{
		id:              uuid.New(),
		trainId:         trainId,
		stationId:       stationId,
		platform:        platform,
		expectedTimes:   exp,
		timetableStatus: status,
	}
}
