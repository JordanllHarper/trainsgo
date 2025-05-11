package main

type train struct {
	entity
	schedule
	name string
}

func newTrain(name string, s station) train {
	return train{
		entity:   newEntity(s.position),
		schedule: schedule{},
		name:     name,
	}
}
