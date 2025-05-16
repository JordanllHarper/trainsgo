package main

import (
	"fmt"
	"maps"
)

type (
	train struct {
		entity
		name string
	}

	trainStoreLocal map[id]train
)

func newTrainStoreLocal() *trainStoreLocal {
	return &trainStoreLocal{}
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

func (tsl trainStoreLocal) changeName(id id, newName string) error {
	train, found := tsl[id]
	if !found {
		return newErrIdNotFound(id, "Train")
	}

	train.name = newName
	tsl[id] = train

	return nil
}

func (t train) String() string {
	return fmt.Sprintf(
		"%v: %v, %v",
		t.name,
		t.entity,
		t.position,
	)
}

func (tsl trainStoreLocal) all() (map[id]train, error) {
	return maps.Clone(tsl), nil

}

func (tsl trainStoreLocal) getByName(name string) ([]train, error) {
	trains := []train{}
	for t := range maps.Values(tsl) {
		if t.name == name {
			trains = append(trains, t)
		}
	}
	return trains, nil
}

func (tsl trainStoreLocal) getById(id id) (train, error) {
	t, found := tsl[id]
	if !found {
		return train{},
			newErrIdNotFound(id, "Train")
	}

	return t, nil
}

func (tsl trainStoreLocal) register(t train) error {
	_, found := tsl[t.id]

	if found {
		return newErrIdAlreadyExists(t.id, "Train")
	}

	tsl[t.id] = t

	return nil
}

func (tsl trainStoreLocal) deregister(id id) error {
	// TODO: Finish this trains schedule and then remove
	return nil
}
