package main

import (
	"encoding/json"
	"log"
	"maps"
	"net/http"
	"os"
	"slices"

	"github.com/google/uuid"
)

type (
	trainHandler struct {
		tReader storeReaderWriter[train]
		sReader storeReader[station]
	}

	stationHandler struct {
		store storeReaderWriter[station]
	}

	lineHandler struct {
		store storeReaderWriter[line]
	}
)

func (th *trainHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	method := req.Method
	log.Printf("Received %s request\n", method)
	switch method {
	case "GET":
		handleGet(rw, req, th.sReader)
	case "POST":
		handleTrainPost(rw, req, th.tReader, th.sReader)
	case "PUT":
		handlePut(rw, req, th.tReader)
	case "DELETE":
		handleDelete(rw, req, th.tReader)
	}
}

func (sh *stationHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	method := req.Method
	log.Printf("Received %s request\n", method)
	switch method {
	case "GET":
		handleGet(rw, req, sh.store)
	case "POST":
		handleStationPost(rw, req, sh.store)
	case "PUT":
		handlePut(rw, req, sh.store)
	case "DELETE":
		handleDelete(rw, req, sh.store)
	}

}

func handleGet[V train | station | line](
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
		js, err := json.Marshal(arr)
		os.Stdout.Write(js)

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

func handleTrainPost(rw http.ResponseWriter, req *http.Request, trw storeReaderWriter[train], ssr storeReader[station]) {
	body := req.Body
	defer body.Close()

	var v struct {
		Name      string `json:"name"`
		StationId string `json:"station_id"`
	}
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

func handleStationPost(rw http.ResponseWriter, req *http.Request, stationStore storeWriter[station]) {
	var t struct {
		Name      string `json:"name"`
		Platforms int    `json:"platforms"`
		X         int    `json:"x"`
		Y         int    `json:"y"`
	}

	err := json.NewDecoder(req.Body).Decode(&t)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	err = stationStore.register(
		newStation(
			newPosition(
				t.X,
				t.Y,
			),
			t.Name,
			t.Platforms,
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

func handlePut[V any](rw http.ResponseWriter, req *http.Request, store storeWriter[V]) {
	var t struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	}

	err := json.NewDecoder(req.Body).Decode(&t)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(t.Id)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	err = store.changeName(id, t.Name)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

func handleDelete[V any](rw http.ResponseWriter, req *http.Request, store storeWriter[V]) {
	var t struct {
		Id string `json:"id"`
	}

	if err := json.NewDecoder(req.Body).Decode(&t); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(t.Id)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	if err = store.deregister(id); err != nil {
		// NOTE: We might end up adding error types depending on the deregister implementations
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)

}
