package main

import (
	"maps"

	"github.com/google/uuid"
)

// Describes a connection between 2 nodes
type (
	Line struct {
		Id         id     `json:"id"`
		Name       string `json:"name"`
		StationOne id     `json:"stationOne"`
		StationTwo id     `json:"stationTwo"`
	}

	lineStoreLocal map[id]Line
)

type lineHandlerLocal struct {
	lines    lineStoreLocal
	stations storeReader[Station]
}

func newLine(one, two Station, name string) Line {
	return Line{
		Id:         uuid.New(),
		StationOne: one.E.Id, StationTwo: two.E.Id,
		Name: name,
	}
}

func newLineStoreLocal() lineStoreLocal {
	return map[id]Line{}
}

func (lsl lineStoreLocal) all() (map[id]Line, *storeReaderError) {
	return maps.Clone(lsl), nil
}

func (lsl lineStoreLocal) getById(id id) (Line, *storeReaderError) {
	value, found := lsl[id]
	if !found {
		return Line{}, newStoreReaderError(id, "Line", StoreReaderErrIdNotFound)
	}
	return value, nil
}

func (lsl lineStoreLocal) getByName(name string) ([]Line, *storeReaderError) {
	lines := []Line{}
	for v := range maps.Values(lsl) {
		if v.Name == name {
			lines = append(lines, v)
		}
	}

	return lines, nil
}

func (lsl lineStoreLocal) changeName(id id, newName string) *storeReaderError {
	line, found := lsl[id]
	if !found {
		return newStoreReaderError(id, "Line", StoreReaderErrIdNotFound)
	}
	line.Name = newName
	lsl[id] = line
	return nil
}

type registerLineErrorCode int

const (
	registerLineErrIdExists registerLineErrorCode = iota
)

type registerLineError struct {
	id   id
	code registerLineErrorCode
}

func (lsl lineStoreLocal) register(l Line) *registerLineError {
	_, found := lsl[l.Id]
	if found {
		return &registerLineError{l.Id, registerLineErrIdExists}
	}

	lsl[l.Id] = l
	return nil
}

func (lsl lineStoreLocal) delete(id id) *storeDeleterError {
	// TODO: Wait for trains to finish using this line, then decommission
	return nil
}
