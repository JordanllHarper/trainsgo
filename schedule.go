package main

import (
	"fmt"
	"maps"

	"github.com/google/uuid"
)

type (
	localScheduler struct {
		schedules map[id]scheduleEntry
	}

	scheduleEntry struct {
		id
		stationId id
		trainId   id
		expectedTimes
	}
)

func newLocalScheduler() *localScheduler {
	return &localScheduler{schedules: map[id]scheduleEntry{}}
}

func newScheduleEntry(stationId id, trainId id, expTimes expectedTimes) scheduleEntry {
	return scheduleEntry{
		id:            uuid.New(),
		stationId:     stationId,
		trainId:       trainId,
		expectedTimes: expTimes,
	}
}

func (sch *localScheduler) add(entry scheduleEntry) error {
	sch.schedules[entry.id] = entry
	return nil
}

func (sch *localScheduler) remove(id id) error {
	delete(sch.schedules, id)

	return nil
}

func (sch *localScheduler) all() (map[id]scheduleEntry, error) {
	return sch.schedules, nil
}

func (sch *localScheduler) getById(id id) (scheduleEntry, error) {
	schedule, found := sch.schedules[id]
	if !found {
		return scheduleEntry{}, newErrIdNotFound(id, "Schedule Entry")
	}

	return schedule, nil
}

func (sch *localScheduler) String() string {
	return fmt.Sprintf("%v", sch.schedules)
}

func (se scheduleEntry) String() string {
	return fmt.Sprintf(
		"station id: %v, trainId: %v, expected arrival: %v, expected departure: %v",
		se.stationId,
		se.trainId,
		se.arrival,
		se.departure,
	)
}
