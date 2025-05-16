package main

import (
	"encoding/json"
	"log"
	"maps"
	"net/http"
	"slices"

	"github.com/google/uuid"
)

type (
	trainHandler struct {
		tReaderWriter storeReaderWriter[Train]
		sReader       storeReader[Station]
	}

	stationHandler struct {
		rw storeReaderWriter[Station]
	}

	navHandler struct {
		lrw storeReaderWriter[Line]
		sr  storeReader[Station]
	}

	trainPostBody struct {
		Name      string `json:"name"`
		StationId string `json:"stationId"`
	}

	stationPostBody struct {
		Name      string `json:"name"`
		Platforms int    `json:"platforms"`
		X         int    `json:"x"`
		Y         int    `json:"y"`
	}

	linePostBody struct {
		Name string `json:"name"`
		One  string `json:"one"`
		Two  string `json:"two"`
	}

	putBody struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	}

	deleteBody struct {
		Id string `json:"id"`
	}
)

func serve[V Train | Station | Line](
	w http.ResponseWriter,
	req *http.Request,
	rw storeReaderWriter[V],
	handlePost func(),
) {
	method := req.Method
	log.Printf("Received %s request\n", method)
	switch method {
	case "GET":
		handleGet(w, req, rw)
	case "POST":
		handlePost()
	case "PUT":
		handlePut(w, req, rw)
	case "DELETE":
		handleDelete(w, req, rw)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func handleGet[V Train | Station | Line](
	rw http.ResponseWriter,
	req *http.Request,
	store storeReader[V],
) {
	query := req.URL.Query()
	hasId := query.Has("id")
	hasName := query.Has("name")

	switch {
	case hasId && hasName:
		http.Error(rw, "Don't include both ID and name", http.StatusBadRequest)
		return

	case !hasId && !hasName:
		all, err := store.all()
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		storeValues := maps.Values(all)
		arr := slices.Collect(storeValues)

		response := struct {
			Values []V `json:"values"`
		}{
			Values: arr,
		}

		if err = json.NewEncoder(rw).Encode(response); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

	case hasId:
		parsedId, err := uuid.Parse(query.Get("id"))
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
		value, err := store.getById(parsedId)

		if err != nil {
			switch err.(type) {
			case errIdNotFound:
				http.Error(rw, err.Error(), http.StatusNotFound)
				return
			default:
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		err = json.NewEncoder(rw).Encode(value)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

	case hasName:
		name := query.Get("name")
		values, err := store.getByName(name)

		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(rw).Encode(values)

		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func handlePut[V any](rw http.ResponseWriter, req *http.Request, store storeWriter[V]) {
	var v putBody

	err := json.NewDecoder(req.Body).Decode(&v)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(v.Id)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	err = store.changeName(id, v.Name)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

func handleDelete[V any](rw http.ResponseWriter, req *http.Request, store storeWriter[V]) {
	var v deleteBody

	if err := json.NewDecoder(req.Body).Decode(&v); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(v.Id)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	if err = store.deregister(id); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)

}

func handleTrainPost(rw http.ResponseWriter, req *http.Request, trw storeReaderWriter[Train], ssr storeReader[Station]) {
	body := req.Body
	defer body.Close()

	var v trainPostBody
	err := json.NewDecoder(body).Decode(&v)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(v.StationId)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	station, err := ssr.getById(id)

	t := newTrain(v.Name, station)

	err = trw.register(t)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusCreated)
}

func handleStationPost(rw http.ResponseWriter, req *http.Request, stationStore storeWriter[Station]) {

	var v stationPostBody

	err := json.NewDecoder(req.Body).Decode(&v)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	err = stationStore.register(
		newStation(
			newPosition(
				v.X,
				v.Y,
			),
			v.Name,
			v.Platforms,
		),
	)

	if err != nil {
		var code int
		switch err.(type) {
		case errStationAlreadyAtPosition, errIdAlreadyExists:
			code = http.StatusBadRequest
		default:
			code = http.StatusInternalServerError
		}
		http.Error(rw, err.Error(), code)
		return
	}
	rw.WriteHeader(http.StatusCreated)
}

func handleNavPost(
	rw http.ResponseWriter,
	req *http.Request,
	lrw storeReaderWriter[Line],
	sr storeReader[Station],
) {
	var v linePostBody

	json.NewDecoder(req.Body).Decode(&v)

	oneId, err := uuid.Parse(v.One)
	twoId, err := uuid.Parse(v.Two)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	st1, err := sr.getById(oneId)
	st2, err := sr.getById(twoId)
	if err != nil {
		switch err.(type) {
		case errIdNotFound:
			http.Error(rw, err.Error(), http.StatusNotFound)
		default:
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	line := newLine(st1, st2, v.Name)
	err = lrw.register(line)
	if err != nil {
		switch err.(type) {
		case errIdAlreadyExists:
			http.Error(rw, err.Error(), http.StatusBadRequest)
		default:
			http.Error(rw, err.Error(), http.StatusInternalServerError)

		}
		return
	}

	rw.WriteHeader(http.StatusCreated)
}
