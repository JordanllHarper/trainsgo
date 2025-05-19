package main

import "fmt"

type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func NewPosition(x, y int) Position {
	return Position{x, y}
}

func (p Position) String() string {
	return fmt.Sprintf("x: %d, y: %d", p.X, p.Y)
}
