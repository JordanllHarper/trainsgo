package main

import (
	"encoding/json"
	"fmt"
	"maps"
	"net/http"

	"github.com/google/uuid"
)

type (
	Train struct {
		E    Entity `json:"entity"`
		Name string `json:"name"`
	}

	trainStoreLocal map[Id]Train
)

func NewTrainStoreLocal() trainStoreLocal {
	return trainStoreLocal{}
}

func NewTrain(
	name string,
	s Station,
) Train {
	return Train{
		E:    NewEntity(s.E.Pos),
		Name: name,
	}
}

func (tsl trainStoreLocal) All() (map[Id]Train, error) {
	return maps.Clone(tsl), nil

}

func (tsl trainStoreLocal) GetById(id Id) (Train, error) {
	t, found := tsl[id]
	if !found {
		return Train{},
			idDoesntExist(id)
	}

	return t, nil
}

func (t Train) String() string {
	return fmt.Sprintf(
		"%v: %v, %v",
		t.Name,
		t.E,
		t.E.Pos,
	)
}

func (tsl trainStoreLocal) register(t Train) error {
	_, found := tsl[t.E.Id]

	if found {
		return idDoesntExist(t.E.Id)
	}

	tsl[t.E.Id] = t

	return nil
}

func (tsl trainStoreLocal) Delete(id Id) error {
	// TODO: Finish this trains schedule and then remove
	return nil
}
func (tsl trainStoreLocal) DeleteBatch(ids []Id) error {
	// TODO: Finish this trains schedule and then remove
	return nil
}

type trainHandlerLocal struct {
	trains   trainStoreLocal
	stations StoreReader[Station]
}

func (h trainHandlerLocal) handlePost(req *http.Request) (HttpResponse, error) {
	body := req.Body
	defer body.Close()

	var v trainPostBody
	err := json.NewDecoder(body).Decode(&v)

	if err != nil {
		return nil, malformedBody{}
	}

	id, err := uuid.Parse(v.StationId)
	if err != nil {
		return nil, badId(v.StationId)
	}
	station, err := h.stations.GetById(id)

	t := NewTrain(v.Name, station)

	err = h.trains.register(t)
	if err != nil {
		return nil, err
	}

	return statusCreated{t}, nil
}
