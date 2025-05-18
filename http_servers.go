package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func (s stationStoreLocal) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	method := req.Method

	var code int
	var body any
	switch method {
	case "GET":
		code, body = handleGet(req, s)
	case "DELETE":
		code, body = handleDelete(req, s)
	case "PUT":
		code, body = s.handlePut(req)
	case "POST":
		code, body = s.handlePost(req)
	default:
		code, body = http.StatusMethodNotAllowed, nil
	}

	serveJson(w, method, code, body)
}

func (h trainHandlerLocal) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	method := req.Method

	var code int
	var body any
	switch method {
	case "GET":
		code, body = handleGet(req, h.trains)
	case "DELETE":
		code, body = handleDelete(req, h.trains)
	case "PUT":
		code, body = h.handlePut(req)
	case "POST":
		code, body = h.handlePost(req)
	default:
		code, body = http.StatusMethodNotAllowed, nil
	}

	serveJson(w, method, code, body)
}

func (h lineHandlerLocal) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	method := req.Method

	var code int
	var body any
	switch method {
	case "GET":
		code, body = handleGet(req, h.lines)
	case "DELETE":
		code, body = handleDelete(req, h.lines)
	case "PUT":
		code, body = h.handlePut(req)
	case "POST":
		code, body = h.handlePost(req)
	default:
		code, body = http.StatusMethodNotAllowed, nil
	}

	serveJson(w, method, code, body)
}

func (h trainHandlerLocal) handlePut(req *http.Request) (int, any) {
	var t renameBody

	if err := json.NewDecoder(req.Body).Decode(&t); err != nil {
		return http.StatusBadRequest, errMalformedBody()
	}

	id, err := uuid.Parse(t.Id)
	if err != nil {
		return http.StatusBadRequest, errBadId(t.Id)
	}

	train, found := h.trains[id]
	if !found {
		return http.StatusBadRequest, errIdDoesntExist(id)
	}

	train.Name = t.Name
	h.trains[id] = train

	return http.StatusOK, train
}

func (h lineHandlerLocal) handlePost(req *http.Request) (int, any) {
	getStationById := func(id string) (Station, any) {
		stId, err := uuid.Parse(id)
		if err != nil {
			return Station{}, errBadId(id)
		}

		station, stErr := h.stations.getById(stId)
		if stErr != nil {
			return Station{}, errIdDoesntExist(stId)
		}
		return station, nil
	}

	var t linePostBody

	if err := json.NewDecoder(req.Body).Decode(&t); err != nil {
		return http.StatusBadRequest, errMalformedBody()
	}

	st1, err := getStationById(t.StationOne)
	if err != nil {
		return http.StatusBadRequest, err
	}

	st2, err := getStationById(t.StationTwo)
	if err != nil {
		return http.StatusBadRequest, err
	}

	line := newLine(st1, st2, t.Name)

	h.lines[line.Id] = line

	return http.StatusOK, line
}

func (h lineHandlerLocal) handlePut(req *http.Request) (int, any) {
	var t renameBody

	if err := json.NewDecoder(req.Body).Decode(&t); err != nil {
		return http.StatusBadRequest, errMalformedBody()
	}

	id, err := uuid.Parse(t.Id)

	if err != nil {
		return http.StatusBadRequest, errBadId(t.Id)
	}

	line, found := h.lines[id]
	if !found {
		return http.StatusBadRequest, errIdDoesntExist(id)
	}

	line.Name = t.Name

	h.lines[id] = line

	return http.StatusOK, line
}

func (s stationStoreLocal) handlePost(req *http.Request) (int, any) {
	var v stationPostBody

	err := json.NewDecoder(req.Body).Decode(&v)
	if err != nil {
		return http.StatusBadRequest, errMalformedBody()
	}

	st :=
		newStation(
			newPosition(
				v.X,
				v.Y,
			),
			v.Name,
			v.Platforms,
		)

	stErr := s.register(st)

	if stErr != nil {
		var errBody errorBody
		switch stErr.code {
		case registerStationErrIdExists:
			// TODO: We shouldn't put this onto the client, we should just try again...
			errBody = errIdExists(stErr.id)
		case registerStationErrPositionTaken:
			errBody = errorBody{Message: "Station Position already taken"}
		default:
			panic(fmt.Sprintf("unexpected main.registerStationErrorCode: %#v", stErr.code))
		}

		return http.StatusBadRequest, errBody
	}
	return http.StatusCreated, st
}
func (s stationStoreLocal) handlePut(req *http.Request) (int, any) {
	var t renameBody

	err := json.NewDecoder(req.Body).Decode(&t)
	if err != nil {
		return http.StatusBadRequest, errMalformedBody()
	}
	id, err := uuid.Parse(t.Id)
	if err != nil {
		return http.StatusBadRequest, errBadId(t.Id)
	}
	value, found := s.stations[id]
	if !found {
		return http.StatusBadRequest, errIdDoesntExist(id)
	}

	value.Name = t.Name
	s.stations[id] = value
	return http.StatusOK, nil
}
