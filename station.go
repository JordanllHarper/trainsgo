package main

type station struct {
	entity
	name             string
	platforms        int
	timetable        map[id]timetableEntry
	surroundingLines []line
}

func newStation(pos position, name string, platforms int) station {
	return station{
		entity:    newEntity(pos),
		name:      name,
		platforms: platforms,
	}
}
