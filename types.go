package main

import (
	"fmt"

	"github.com/google/uuid"
)

type (
	id = uuid.UUID

	position struct {
		x, y int
	}

	entity struct {
		id
		position
	}
)

func newEntity(pos position) entity {
	return entity{
		id:       uuid.New(),
		position: pos,
	}
}

func (p position) String() string {
	return fmt.Sprintf("x: %d, y: %d", p.x, p.y)
}
func (e entity) String() string {
	return fmt.Sprintf("Id: %v - %v", e.id, e.position)
}
