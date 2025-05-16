package main

import (
	"fmt"

	"github.com/google/uuid"
)

type entity struct {
	id       `json:"id"`
	position `json:"position"`
}

func newEntity(pos position) entity {
	return entity{
		id:       uuid.New(),
		position: pos,
	}
}

func (e entity) String() string {
	return fmt.Sprintf("Id: %v - %v", e.id, e.position)
}
