package engine

import (
	"testing"
)

func TestBuildGrid(t *testing.T) {
	expectedLenX := 1000
	expectedLenY := 1000
	actual := BuildCellGrid(map[Coordinates]bool{})

	testY, testX := len(actual), len(actual[0])

	if testY != expectedLenY {
		t.Fatalf("Invalid Y length, wanted %v, got %v", expectedLenY, testY)
		return
	}

	if testX != expectedLenX {
		t.Fatalf("Invalid X length, wanted %v, got %v", expectedLenX, testX)
		return
	}
}
