package common

import (
	"fmt"
	"strconv"
)

type Status string

const (
	Travelling   Status = "Travelling"
	Transferring        = "Transferring"
	Unused              = "Unused"
	Emergency           = "Emergency"
)

type Train struct {
	Name string `json:"name"`
	// meters / second
	TopSpeed int `json:"top_speed"`
	Coordinates

	Status `json:"status"`
}

func NewTrain(name string, topSpeed int, coords Coordinates, status Status) Train {
	return Train{
		name,
		topSpeed,
		coords,
		status,
	}
}

func (t Train) ToString() string {
	return fmt.Sprintf("Name: %s\nTop speed: %s\nCoords: %s\nStatus: %s", t.Name, strconv.Itoa(t.TopSpeed), t.Coordinates.Pretty(), t.Status)
}
