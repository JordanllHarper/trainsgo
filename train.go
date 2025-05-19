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

func newTrainStoreLocal() trainStoreLocal {
	return trainStoreLocal{}
}

func newTrain(
	name string,
	s Station,
) Train {
	return Train{
		E:    NewEntity(s.E.Pos),
		Name: name,
	}
}

func (tsl trainStoreLocal) all() (map[Id]Train, *storeReaderError) {
	return maps.Clone(tsl), nil

}

func (tsl trainStoreLocal) getById(id Id) (Train, *storeReaderError) {
	t, found := tsl[id]
	if !found {
		return Train{},
			newStoreReaderError(id, "Train", StoreReaderErrIdNotFound)
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

type registerTrainErrorCode int

const (
	registerTrainErrIdExists registerStationErrorCode = iota
)

type registerTrainError struct {
	id   Id
	code registerStationErrorCode
}

func (tsl trainStoreLocal) register(t Train) *registerTrainError {
	_, found := tsl[t.E.Id]

	if found {
		return &registerTrainError{t.E.Id, registerTrainErrIdExists}
	}

	tsl[t.E.Id] = t

	return nil
}

func (tsl trainStoreLocal) delete(id Id) *storeDeleterError {
	// TODO: Finish this trains schedule and then remove
	return nil
}

type trainHandlerLocal struct {
	trains   trainStoreLocal
	stations storeReader[Station]
}

func (h trainHandlerLocal) handlePost(req *http.Request) (int, any) {
	body := req.Body
	defer body.Close()

	var v trainPostBody
	err := json.NewDecoder(body).Decode(&v)

	if err != nil {
		return http.StatusBadRequest, errMalformedBody()
	}

	id, err := uuid.Parse(v.StationId)
	if err != nil {
		return http.StatusBadRequest, errorBody{Message: fmt.Sprintf("Bad Station ID: %s", id)}
	}
	station, err := h.stations.getById(id)

	t := newTrain(v.Name, station)

	trErr := h.trains.register(t)
	if trErr != nil {
		var errBody errorBody
		switch trErr.code {
		case registerTrainErrIdExists:
			errIdExists(trErr.id)
		default:
			panic(fmt.Sprintf("unexpected main.registerStationErrorCode: %#v", trErr.code))
		}
		return http.StatusInternalServerError, errBody
	}

	return http.StatusCreated, nil
}
