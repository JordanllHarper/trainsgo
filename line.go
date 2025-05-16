package main

import (
	"maps"

	"github.com/google/uuid"
)

type (
	// Describes a connection between 2 nodes
	line struct {
		Id   id     `json:"id"`
		Name string `json:"name"`
		One  Node   `json:"one"`
		Two  Node   `json:"two"`
	}

	// Describes a point where multiple connections can interact
	intersection struct {
		entity
		connections map[id]line
	}

	navigationStoreLocal struct {
		lines         map[id]line
		intersections map[id]intersection
	}
)

func newNavigationStoreLocal() *navigationStoreLocal {
	return &navigationStoreLocal{
		lines:         map[id]line{},
		intersections: map[id]intersection{},
	}
}

func newLine(one, two Node, name string) line {
	return line{
		Id:  uuid.New(),
		One: one, Two: two,
		Name: name,
	}
}

func newIntersection(pos position, connections map[id]line) intersection {
	return intersection{
		entity:      newEntity(pos),
		connections: connections,
	}
}

func (nsl *navigationStoreLocal) all() (map[id]line, error) {
	return maps.Clone(nsl.lines), nil
}

func (nsl *navigationStoreLocal) getById(id id) (line, error) {
	value, found := nsl.lines[id]
	if !found {
		return line{}, newErrIdNotFound(id, "Line")
	}
	return value, nil
}

func (nsl *navigationStoreLocal) getByName(name string) ([]line, error) {
	lines := []line{}
	for v := range maps.Values(nsl.lines) {
		if v.Name == name {
			lines = append(lines, v)
		}
	}

	return lines, nil
}

func (nsl *navigationStoreLocal) register(l line) error {
	_, found := nsl.lines[l.Id]
	if found {
		return newErrIdAlreadyExists(l.Id, "Line")
	}

	nsl.lines[l.Id] = l
	return nil
}

func (nsl *navigationStoreLocal) deregister(id id) error {
	// TODO: Wait for trains to finish using this line, then decommission
	return nil
}
