package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func (s *stationStoreLocal) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	method := req.Method

	var response HttpResponse
	var err error
	switch method {
	case "GET":
		response, err = handleGet(req, s)
	case "DELETE":
		response, err = handleDelete(req, s)
	case "PUT":
		response, err = s.handlePut(req)
	case "POST":
		response, err = s.handlePost(req)
	default:
		response, err = nil, methodNotAllowed(method)
	}

	serveJson(w, method, response, err)
}

func (h trainHandlerLocal) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	method := req.Method

	var code HttpResponse
	var body error
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
		code, body = nil, methodNotAllowed(method)
	}

	serveJson(w, method, code, body)
}

func (h lineHandlerLocal) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	method := req.Method

	var code HttpResponse
	var body error
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
		code, body = nil, methodNotAllowed(method)
	}

	serveJson(w, method, code, body)
}

func (h tripHandlerLocal) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	method := req.Method

	var code HttpResponse
	var body error
	switch method {
	case "GET":
		code, body = handleGet(req, h.trips)
	case "POST":
		code, body = h.handlePost(req)
	case "PUT":
		code, body = h.handlePut(req)
	}

	serveJson(w, method, code, body)
}

func (h trainHandlerLocal) handlePut(req *http.Request) (HttpResponse, HttpError) {
	var t renameBody

	if err := json.NewDecoder(req.Body).Decode(&t); err != nil {
		return nil, malformedBody{}
	}

	id, err := uuid.Parse(t.Id)
	if err != nil {
		return nil, badId(t.Id)
	}

	train, found := h.trains[id]
	if !found {
		return nil, idDoesntExist(id)
	}

	train.Name = t.Name
	h.trains[id] = train

	return statusOK{body: train}, nil
}

func (h lineHandlerLocal) handlePost(req *http.Request) (HttpResponse, error) {

	var t linePostBody

	if err := json.NewDecoder(req.Body).Decode(&t); err != nil {
		return nil, malformedBody{}
	}

	st1, err := getByStringId(h.stations, t.StationOne)
	if err != nil {
		return nil, err
	}

	st2, err := getByStringId(h.stations, t.StationTwo)
	if err != nil {
		return nil, err
	}

	line := newLine(st1, st2, t.Name)

	h.lines[line.Id] = line

	return statusOK{line}, nil
}

func (h lineHandlerLocal) handlePut(req *http.Request) (HttpResponse, HttpError) {
	var t renameBody

	if err := json.NewDecoder(req.Body).Decode(&t); err != nil {
		return nil, malformedBody{}
	}

	id, err := uuid.Parse(t.Id)

	if err != nil {
		return nil, badId(t.Id)
	}

	line, found := h.lines[id]
	if !found {
		return nil, idDoesntExist(id)
	}

	line.Name = t.Name

	h.lines[id] = line

	return statusOK{line}, nil
}

func (s stationStoreLocal) handlePost(req *http.Request) (HttpResponse, error) {
	var v stationPostBody

	err := json.NewDecoder(req.Body).Decode(&v)
	if err != nil {
		return nil, malformedBody{}
	}

	st :=
		NewStation(
			NewPosition(
				v.X,
				v.Y,
			),
			v.Name,
			v.Platforms,
		)

	err = s.register(st)

	if err != nil {
		return nil, err

	}
	return statusCreated{st}, nil
}

func (s stationStoreLocal) handlePut(req *http.Request) (HttpResponse, HttpError) {
	var t renameBody

	err := json.NewDecoder(req.Body).Decode(&t)
	if err != nil {
		return nil, malformedBody{}
	}
	id, err := uuid.Parse(t.Id)
	if err != nil {
		return nil, badId(t.Id)
	}
	value, found := s.stations[id]
	if !found {
		return nil, idDoesntExist(id)
	}

	value.Name = t.Name
	s.stations[id] = value
	return statusOK{}, nil
}

func (h tripHandlerLocal) handlePost(req *http.Request) (HttpResponse, error) {
	var t tripPostBody

	if err := json.NewDecoder(req.Body).Decode(&t); err != nil {
		return nil, malformedBody{}
	}

	fromStation, err := getByStringId(h.stations, t.FromStationId)
	if err != nil {
		return nil, err
	}

	toStation, err := getByStringId(h.stations, t.ToStationId)
	if err != nil {
		return nil, err
	}

	train, err := getByStringId(h.trains, t.TrainId)
	if err != nil {
		return nil, err
	}

	trip := NewTrip(fromStation, toStation, train, t.ExpTimes, t.StartingStatus)

	h.trips[trip.Id] = trip

	if _, err := h.router.MapRoute(fromStation.E.Id, toStation.E.Id); err != nil {
		return nil, err
	}
	// We've added trip to the store, now we need to plot a route from the starting station to the ending station
	return statusCreated{trip}, nil
}

func (h tripHandlerLocal) handlePut(req *http.Request) (HttpResponse, error) {
	var t tripPutBody

	if err := json.NewDecoder(req.Body).Decode(&t); err != nil {
		return nil, malformedBody{}
	}

	trip, err := getByStringId(h.trips, t.Id)
	if err != nil {
		return nil, err
	}

	trip.Status = t.NewStatus
	h.trips[trip.Id] = trip

	return statusOK{trip}, nil

}
