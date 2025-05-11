package main

import "github.com/google/uuid"

type (
	// Describes a connection between 2 nodes
	line struct {
		id
		one, two node
	}

	intersection struct {
		entity
		connections []line
	}
)

func newLine(one, two node) line {
	return line{
		id:  uuid.New(),
		one: one, two: two,
	}
}

func newIntersection(pos position, connections []line) intersection {
	return intersection{
		entity:      newEntity(pos),
		connections: connections,
	}
}
