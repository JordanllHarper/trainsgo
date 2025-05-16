package main

import (
	"fmt"
	"maps"
	"time"

	"github.com/google/uuid"
)

type (
	timetableViewer interface {
		all() ([]timetableEntry, error)
		getById(id id) (timetableEntry, error)
	}

	timetableHandler interface {
		add(e timetableEntry) error
		delay(id id, by time.Duration) error
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

	timetableHandlerImpl struct {
		timetable map[id]timetableEntry
	}
)

func newTimetableHandlerImpl() *timetableHandlerImpl {
	return &timetableHandlerImpl{timetable: map[id]timetableEntry{}}
}

/*
Delays a timetable by a given duration and returns it.

Assumes that the arrival time isSome.
*/
func (entry timetableEntry) delayed(by time.Duration) timetableEntry {
	newEntryTime := entry.arrival.value.Add(by)
	entry.timetableStatus = Delayed
	entry.arrival = newSomeTime(newEntryTime)
	return entry
}

func (t timetableHandlerImpl) all() (map[id]timetableEntry, error) {
	// When viewing, we don't want to expose an underlying slice to the caller, this should be "by value" slice (they get a copy)
	return maps.Clone(t.timetable), nil
}

func (t timetableHandlerImpl) getById(id id) (timetableEntry, error) {

	value, found := t.timetable[id]

	if !found {
		return timetableEntry{}, newErrIdNotFound(id, "Timetable entry")
	}

	return value, nil
}

func (t timetableHandlerImpl) add(e timetableEntry) error {

	_, found := t.timetable[e.id]

	if found {
		return newErrIdAlreadyExists(e.id, "Timetable Entry")
	}

	/*
		TODO: Prevent edge cases:
		- if the platform and time range is already in use, conflict!
		- think about this more ...
	*/

	t.timetable[e.id] = e

	return nil
}

type delayWithNoArrivalTime struct {
	id
	delay time.Duration
}

func (err delayWithNoArrivalTime) Error() string {
	return fmt.Sprintf("Attempted to delay Timetable Entry ID %v by %v, but arrivalTime was None", err.id, err.delay)
}

func newDelayWithNoArrivalTime(id id, delay time.Duration) error {
	return delayWithNoArrivalTime{id, delay}
}

func (t timetableHandlerImpl) delay(id id, by time.Duration) error {

	value, found := t.timetable[id]

	if !found {
		return newErrIdNotFound(id, "Timetable Entry")
	}

	if !value.arrival.isSome {
		// we can't delay if the entry doesn't have an arrival time
		return newDelayWithNoArrivalTime(id, by)
	}

	t.timetable[id] = value.delayed(by)

	return nil
}

func (t timetableHandlerImpl) cancelAndRemove(id id) error {
	value, found := t.timetable[id]

	value.timetableStatus = Cancelled

	if !found {
		return newErrIdNotFound(id, "Timetable Entry")
	}

	t.timetable[id] = value

	return nil
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
