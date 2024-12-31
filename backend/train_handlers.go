package main

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
)

func needBody() ClientError      { return ClientError{"Need body"} }
func provideId() ClientError     { return ClientError{"No ID provided"} }
func invalidId() ClientError     { return ClientError{"Invalid ID"} }
func malformedBody() ClientError { return ClientError{"Malformed body"} }

type TrainResponseEmpty struct {
	int
}

func (response TrainResponseEmpty) StatusCode() int {
	return response.int
}

type TrainResponseSingular struct {
	TrainEntity
	statusCode int
}

func (response TrainResponseSingular) StatusCode() int {
	return response.statusCode
}

type TrainResponseMultiple struct {
	Trains     []TrainEntity
	statusCode int
}

func (response TrainResponseMultiple) StatusCode() int {
	return response.statusCode
}

type TrainGetRequest interface{}

// Gets a train.
//
// Accepts an "id" in a request for a specified train, or leave blank for all the trains available.
func onTrainGet(db *gorm.DB, req *http.Request) (ResponseBody, HttpError) {
	queries := req.URL.Query()
	id := queries.Get("id")
	if strings.TrimSpace(id) == "" {
		var trainEntities []TrainEntity
		db.Find(&trainEntities)
		return TrainResponseMultiple{Trains: trainEntities, statusCode: http.StatusOK}, nil
	}
	parsedId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, provideId()
	}

	var trainEntity TrainEntity
	result := db.Where(&TrainEntity{DbFields: DbFields{ID: uint(parsedId)}}).Find(&trainEntity)

	if result.RowsAffected == 0 {
		return TrainResponseEmpty{http.StatusNoContent}, nil
	}

	return TrainResponseSingular{trainEntity, http.StatusOK}, nil
}

// Creates a new train and inserts into the database. Returns the train in the body of the response.
func onTrainPost(db *gorm.DB, req *http.Request) (ResponseBody, HttpError) {
	body := req.Body
	if body == nil {
		return nil, needBody()
	}

	var train Train
	decoder := json.NewDecoder(body)
	err := decoder.Decode(&train)
	if err != nil {
		return nil, malformedBody()
	}

	tEntity := TrainEntity{DbFields: DbFields{}, Train: train}
	db = db.Create(&tEntity)
	fmt.Printf("[INFO]: Inserted train entity: %v\n", tEntity)

	return TrainResponseSingular{tEntity, http.StatusCreated}, nil
}

func onTrainPut(db *gorm.DB, req *http.Request) (ResponseBody, HttpError) {

	body := req.Body
	queries := req.URL.Query()
	id := queries.Get("id")
	if strings.TrimSpace(id) == "" {
		return nil, provideId()
	}

	parsedId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, invalidId()
	}

	var updatedTrainFields Train
	decoder := json.NewDecoder(body)
	err = decoder.Decode(&updatedTrainFields)
	if err != nil {
		return nil, malformedBody()
	}

	trainEntity := TrainEntity{DbFields: DbFields{ID: uint(parsedId)}, Train: updatedTrainFields}
	result := db.Updates(&trainEntity)
	if result.RowsAffected != 0 {
		fmt.Printf("[INFO]: Modified train entity id: %v\n", parsedId)
		return TrainResponseSingular{trainEntity, http.StatusOK}, nil
	} else {
		fmt.Printf("[INFO]: Attempted modification with no result: %v\n", parsedId)
		return TrainResponseEmpty{http.StatusNoContent}, nil
	}
}

func onTrainDelete(db *gorm.DB, req *http.Request) (ResponseBody, HttpError) {
	queries := req.URL.Query()
	id := queries.Get("id")
	if strings.TrimSpace(id) == "" {
		return nil, provideId()
	}

	parsedId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, invalidId()
	}

	trainEntity := &TrainEntity{
		DbFields: DbFields{ID: uint(parsedId)},
	}
	result := db.Delete(&trainEntity)
	if result.RowsAffected != 0 {
		fmt.Printf("[INFO]: Deleted train entity id: %v\n", parsedId)
	} else {
		fmt.Printf("[INFO]: Attempted deletion with no result: %v\n", parsedId)
	}

	return TrainResponseEmpty{http.StatusNoContent}, nil
}
