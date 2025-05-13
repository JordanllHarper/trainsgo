package main

import "fmt"

type train struct {
	entity
	name string
}

func newTrain(
	name string,
	s station,
) train {
	return train{
		entity: newEntity(s.position),
		name:   name,
	}
}

func (t train) String() string {
	return fmt.Sprintf(
		"%v: %v, %v",
		t.name,
		t.entity,
		t.position,
	)
}
