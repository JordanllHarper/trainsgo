package main

import "fmt"

type (
	scheduleViewer interface {
		all() ([]scheduleEntry, error)
	}

	scheduler interface {
		add(s scheduleEntry) error
		remove() (scheduleEntry, error)
	}

	localScheduler struct {
		schedules []scheduleEntry
	}

	scheduleEntry struct {
		stationId id
		trainId   id
		expectedTimes
	}
)

func newScheduleEntry(stationId id, trainId id, expTimes expectedTimes) scheduleEntry {
	return scheduleEntry{stationId, trainId, expTimes}
}

func (sch *localScheduler) add(entry scheduleEntry) error {
	newQueue := append(sch.schedules, entry)
	sch.schedules = newQueue
	return nil
}

func (sch *localScheduler) remove() (scheduleEntry, error) {
	queue := sch.schedules
	deref := queue
	entry := deref[0]
	newQueue := deref[1:]
	sch.schedules = newQueue
	return entry, nil
}

func (sch *localScheduler) all() ([]scheduleEntry, error) {
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
