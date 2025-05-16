package main

import (
	"fmt"
	"maps"
)

type (
	train struct {
		E    entity `json:"entity"`
		Name string `json:"name"`
	}

	trainStoreLocal map[id]train
)

func newTrainStoreLocal() *trainStoreLocal {
	return &trainStoreLocal{}
}

func newTrain(
	name string,
	s Station,
) train {
	return train{
		E:    newEntity(s.E.Pos),
		Name: name,
	}
}

func (tsl trainStoreLocal) changeName(id id, newName string) error {
	train, found := tsl[id]
	if !found {
		return newErrIdNotFound(id, "Train")
	}

	train.Name = newName
	tsl[id] = train

	return nil
}

func (t train) String() string {
	return fmt.Sprintf(
		"%v: %v, %v",
		t.Name,
		t.E,
		t.E.Pos,
	)
}

func (tsl trainStoreLocal) all() (map[id]train, error) {
	return maps.Clone(tsl), nil

}

func (tsl trainStoreLocal) getByName(name string) ([]train, error) {
	trains := []train{}
	for t := range maps.Values(tsl) {
		if t.Name == name {
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
	_, found := tsl[t.E.Id]

	if found {
		return newErrIdAlreadyExists(t.E.Id, "Train")
	}

	tsl[t.E.Id] = t

	return nil
}

func (tsl trainStoreLocal) deregister(id id) error {
	// TODO: Finish this trains schedule and then remove
	return nil
}
