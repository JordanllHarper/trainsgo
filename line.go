package main

import (
	"maps"

	"github.com/google/uuid"
)

// Describes a connection between 2 nodes
type (
	Line struct {
		Id         Id     `json:"id"`
		Name       string `json:"name"`
		StationOne Id     `json:"stationOne"`
		StationTwo Id     `json:"stationTwo"`
	}

	lineStoreLocal map[Id]Line
)

type lineHandlerLocal struct {
	lines    lineStoreLocal
	stations StoreReader[Station]
}

func newLine(one, two Station, name string) Line {
	return Line{
		Id:         uuid.New(),
		StationOne: one.E.Id, StationTwo: two.E.Id,
		Name: name,
	}
}

func NewLineStoreLocal() lineStoreLocal {
	return map[Id]Line{}
}

func (lsl lineStoreLocal) All() (map[Id]Line, *StoreReaderError) {
	return maps.Clone(lsl), nil
}

func (lsl lineStoreLocal) GetById(id Id) (Line, *StoreReaderError) {
	value, found := lsl[id]
	if !found {
		return Line{}, NewStoreReaderError(id, "Line", StoreReaderErrIdNotFound)
	}
	return value, nil
}

func (lsl lineStoreLocal) getByName(name string) ([]Line, *StoreReaderError) {
	lines := []Line{}
	for v := range maps.Values(lsl) {
		if v.Name == name {
			lines = append(lines, v)
		}
	}

	return lines, nil
}

func (lsl lineStoreLocal) changeName(id Id, newName string) *StoreReaderError {
	line, found := lsl[id]
	if !found {
		return NewStoreReaderError(id, "Line", StoreReaderErrIdNotFound)
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
	id   Id
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

func (lsl lineStoreLocal) Delete(id Id) *StoreDeleterError {
	// TODO: Wait for trains to finish using this line, then decommission
	return nil
}
