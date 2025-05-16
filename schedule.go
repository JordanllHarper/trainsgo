package main

import (
	"fmt"

	"github.com/google/uuid"
)

type (
	scheduleViewer interface {
		all() (map[id]scheduleEntry, error)
	}

	scheduler interface {
		add(s scheduleEntry) error
		remove(id) error
	}

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
