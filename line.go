package main

import (
	"maps"
	"strings"

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
	stations StoreIDable[Station]
}

func newLine(one, two Station, name string) Line {
	return Line{
		Id:         uuid.New(),
		StationOne: one.E.Id,
		StationTwo: two.E.Id,
		Name:       name,
	}
}

func NewLineStoreLocal() lineStoreLocal {
	return map[Id]Line{}
}

func (lsl lineStoreLocal) All() (map[Id]Line, error) {
	return maps.Clone(lsl), nil
}

func (lsl lineStoreLocal) GetById(id Id) (Line, error) {
	value, found := lsl[id]
	if !found {
		return Line{}, idDoesntExist(id)
	}
	return value, nil
}

func (lsl lineStoreLocal) GetByName(name string) ([]Line, error) {
	lines := []Line{}
	for v := range maps.Values(lsl) {
		if strings.Contains(v.Name, name) {
			lines = append(lines, v)
		}
	}

	return lines, nil
}

func (lsl lineStoreLocal) changeName(id Id, newName string) error {
	line, found := lsl[id]
	if !found {
		return idDoesntExist(id)
	}
	line.Name = newName
	lsl[id] = line
	return nil
}

func (lsl lineStoreLocal) register(l Line) error {
	_, found := lsl[l.Id]
	if found {
		return idAlreadyExists(l.Id)
	}

	lsl[l.Id] = l
	return nil
}

func (lsl lineStoreLocal) Delete(id Id) error {
	// TODO: Wait for trains to finish using this line, then decommission
	return nil
}

func (lsl lineStoreLocal) DeleteBatch(ids []Id) error {
	// TODO: Wait for trains to finish using these lines, then decommission
	return nil
}
