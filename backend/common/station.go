package common

type Station struct {
	Name string `json:"name"`
	Coordinates
}

func NewStation(name string, coords Coordinates) Station {
	return Station{name, coords}
}
