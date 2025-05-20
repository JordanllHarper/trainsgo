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

func (tsl trainStoreLocal) All() (map[Id]Train, StoreError) {
	return maps.Clone(tsl), nil

}

func (tsl trainStoreLocal) GetById(id Id) (Train, StoreError) {
	t, found := tsl[id]
	if !found {
		return Train{},
			IdDoesntExist(id)
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

type RegisterTrainErrorCode int

const (
	RegisterTrainErrIdAlreadyExist RegisterTrainErrorCode = 0
)

type RegisterTrainError interface {
	error
	RegisterCode() RegisterTrainErrorCode
}

func (e IdDoesntExist) Id() Id {
	return Id(e)
}

func (e IdDoesntExist) RegisterCode() RegisterTrainErrorCode {
	return RegisterTrainErrIdAlreadyExist
}

func (tsl trainStoreLocal) register(t Train) RegisterTrainError {
	_, found := tsl[t.E.Id]

	if found {
		return IdDoesntExist(t.E.Id)
	}

	tsl[t.E.Id] = t

	return nil
}

func (tsl trainStoreLocal) Delete(id Id) StoreError {
	// TODO: Finish this trains schedule and then remove
	return nil
}

type trainHandlerLocal struct {
	trains   trainStoreLocal
	stations StoreReader[Station]
}

func (h trainHandlerLocal) handlePost(req *http.Request) (HttpResponse, HttpError) {
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

	trErr := h.trains.register(t)
	if trErr != nil {
		switch trErr.RegisterCode() {
		case RegisterTrainErrIdAlreadyExist:
			return nil, idAlreadyExists(t.E.Id)
		}
		return nil, internalServerError{trErr}
	}

	return statusCreated{t}, nil
}
