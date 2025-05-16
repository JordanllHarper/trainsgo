package main

import (
	"fmt"

	"github.com/google/uuid"
)

type entity struct {
	Id  id       `json:"id"`
	Pos position `json:"pos"`
}

func newEntity(pos position) entity {
	return entity{
		Id:  uuid.New(),
		Pos: pos,
	}
}

func (e entity) String() string {
	return fmt.Sprintf("Id: %v - %v", e.Id, e.Pos)
}
