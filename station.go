package main

import "fmt"

type station struct {
	entity
	name             string
	platforms        int
	surroundingLines []line
}

func newStation(pos position, name string, platforms int) station {
	return station{
		entity:    newEntity(pos),
		name:      name,
		platforms: platforms,
	}
}

func (s station) String() string {
	return fmt.Sprintf("%v: num_platforms - %v", s.name, s.platforms)
}
