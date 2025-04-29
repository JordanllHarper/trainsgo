package common

import (
	"fmt"
	"strconv"
)

// Coordinates for locating on an x and y axis
type Coordinates struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func (coord Coordinates) Pretty() string {
	return fmt.Sprintf("%s, %s", strconv.Itoa(coord.X), strconv.Itoa(coord.Y))
}
