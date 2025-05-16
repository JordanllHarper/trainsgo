package main

import (
	"maps"

	"github.com/google/uuid"
)

// Describes a connection between 2 nodes
type (
	Line struct {
		Id   id     `json:"id"`
		Name string `json:"name"`
		One  id     `json:"one"`
		Two  id     `json:"two"`
	}

	lineStoreLocal struct {
		lines map[id]Line
	}
)

func newLine(one, two Station, name string) Line {
	return Line{
		Id:  uuid.New(),
		One: one.E.Id, Two: two.E.Id,
		Name: name,
	}
}

func newLineStoreLocal() *lineStoreLocal {
	return &lineStoreLocal{
		lines: map[id]Line{},
	}
}

func (nsl *lineStoreLocal) all() (map[id]Line, error) {
	return maps.Clone(nsl.lines), nil
}

func (nsl *lineStoreLocal) getById(id id) (Line, error) {
	value, found := nsl.lines[id]
	if !found {
		return Line{}, newErrIdNotFound(id, "Line")
	}
	return value, nil
}

func (nsl *lineStoreLocal) getByName(name string) ([]Line, error) {
	lines := []Line{}
	for v := range maps.Values(nsl.lines) {
		if v.Name == name {
			lines = append(lines, v)
		}
	}

	return lines, nil
}

func (lsl *lineStoreLocal) changeName(id id, newName string) error {
	line, found := lsl.lines[id]
	if !found {
		return newErrIdNotFound(id, "Line")
	}
	line.Name = newName
	lsl.lines[id] = line
	return nil
}

func (nsl *lineStoreLocal) register(l Line) error {
	_, found := nsl.lines[l.Id]
	if found {
		return newErrIdAlreadyExists(l.Id, "Line")
	}

	nsl.lines[l.Id] = l
	return nil
}

func (nsl *lineStoreLocal) deregister(id id) error {
	// TODO: Wait for trains to finish using this line, then decommission
	return nil
}
