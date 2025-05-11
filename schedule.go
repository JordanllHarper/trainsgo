package main

type (
	schedule struct {
		scheduleQueue []scheduleEntry
	}

	scheduleEntry struct {
		stationId id
		expectedTimes
	}
)

func newScheduleEntry(stationId id, expTimes expectedTimes) scheduleEntry {
	return scheduleEntry{stationId, expTimes}
}

func (sch *schedule) push(entry scheduleEntry) {
	newQueue := append(sch.scheduleQueue, entry)
	sch.scheduleQueue = newQueue
}

func (sch *schedule) pop() scheduleEntry {
	queue := sch.scheduleQueue
	deref := queue
	entry := deref[0]
	newQueue := deref[1:]
	sch.scheduleQueue = newQueue
	return entry
}

func (sch *schedule) view() []scheduleEntry {
	return sch.scheduleQueue
}
