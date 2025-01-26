package engine

import ()

/*
TODO:

	Setup a channel that allows creation of stations on said grid
	Setup a channel that allows for creation of trains between stations.
	Simulate train moving across grid - hard bit
*/

/*
NOTE:
Grid is not necessary, we just need to keep track of entities. We don't need to actually plot them on a grid.

*/

const CellGridSize = 1000

type CellType int
type Information string

const (
	Empty CellType = iota
	Station
)

type Cell struct {
	CellType
	Information
}

func NewEmptyCell() Cell   { return Cell{CellType: Empty, Information: "Empty Cell"} }
func NewStationCell() Cell { return Cell{CellType: Station, Information: "Station Cell"} }
