package main

import (
	"fmt"

	"github.com/google/uuid"
)

type Entity struct {
	Id  Id       `json:"id"`
	Pos Position `json:"pos"`
}

func NewEntity(pos Position) Entity {
	return Entity{
		Id:  uuid.New(),
		Pos: pos,
	}
}

func (e Entity) String() string {
	return fmt.Sprintf("Id: %v - %v", e.Id, e.Pos)
}
