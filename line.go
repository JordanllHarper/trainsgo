package main

import (
	"maps"

	"github.com/google/uuid"
)

type (
	// Describes a connection between 2 nodes
	Line struct {
		Id   id     `json:"id"`
		Name string `json:"name"`
		One  Node   `json:"one"`
		Two  Node   `json:"two"`
	}

	// Describes a point where multiple connections can interact
	intersection struct {
		E           entity      `json:"entity"`
		Connections map[id]Line `json:"connections"`
	}

	navigationStoreLocal struct {
		lines         map[id]Line
		intersections map[id]intersection
	}
)

func newNavigationStoreLocal() *navigationStoreLocal {
	return &navigationStoreLocal{
		lines:         map[id]Line{},
		intersections: map[id]intersection{},
	}
}

func newLine(one, two Node, name string) Line {
	return Line{
		Id:  uuid.New(),
		One: one, Two: two,
		Name: name,
	}
}

func newIntersection(pos position, connections map[id]Line) intersection {
	return intersection{
		E:           newEntity(pos),
		Connections: connections,
	}
}

func (nsl *navigationStoreLocal) all() (map[id]Line, error) {
	return maps.Clone(nsl.lines), nil
}

func (nsl *navigationStoreLocal) getById(id id) (Line, error) {
	value, found := nsl.lines[id]
	if !found {
		return Line{}, newErrIdNotFound(id, "Line")
	}
	return value, nil
}

func (nsl *navigationStoreLocal) getByName(name string) ([]Line, error) {
	lines := []Line{}
	for v := range maps.Values(nsl.lines) {
		if v.Name == name {
			lines = append(lines, v)
		}
	}

	return lines, nil
}

func (nsl *navigationStoreLocal) register(l Line) error {
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
