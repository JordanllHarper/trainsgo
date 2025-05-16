package main

import "fmt"

type position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func newPosition(x, y int) position {
	return position{x, y}
}

func (p position) String() string {
	return fmt.Sprintf("x: %d, y: %d", p.X, p.Y)
}
