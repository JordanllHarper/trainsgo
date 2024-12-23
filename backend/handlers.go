package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

// Root requests provide information about the state of the applicationa
// Used in the train dashboard.
func handleRoot(w http.ResponseWriter, r *http.Request) {}

type TrainGetRequestEmpty struct{}

type TrainGetRequestSingular struct {
	TrainEntity
}

type TrainGetRequestMultiple struct {
	Trains []TrainEntity
}

type TrainGetRequest interface{}

// Gets a train.
//
// Accepts an "id" in a request for a specified train, or leave blank for all the trains available.
func onTrainGet(db *gorm.DB, req *http.Request) (ResponseBody, HttpError) {
	queries := req.URL.Query()
	id := queries.Get("id")
	if id == "" {
		var trainEntities []TrainEntity
		db.Find(&trainEntities)
		return TrainGetRequestMultiple{trainEntities}, nil
	}
	parsedId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, ClientError{"Invalid Id"}
	}

	var trainEntity TrainEntity
	result := db.Where(&TrainEntity{DbFields: DbFields{ID: uint(parsedId)}}).Find(&trainEntity)

	if result.RowsAffected == 0 {
		return TrainGetRequestEmpty{}, nil
	}

	return TrainGetRequestSingular{trainEntity}, nil
}

var needBody = ClientError{"Need body"}
var malformedBody = ClientError{"Malformed body"}

// Creates a new train and inserts into the database. Returns the train in the body of the response.
func onTrainPost(db *gorm.DB, req *http.Request) (ResponseBody, HttpError) {
	body := req.Body

	if body == nil {
		return nil, needBody
	}

	decoder := json.NewDecoder(body)

	var train Train
	err := decoder.Decode(&train)

	if err != nil {
		return nil, malformedBody
	}

	tEntity := TrainEntity{DbFields: DbFields{}, Train: train}
	db.Create(&tEntity)
	fmt.Printf("[INFO]: Inserted train entity: %v\n", tEntity)

	return TrainGetRequestSingular{tEntity}, nil
}

var provideId ClientError = ClientError{"No ID provided"}
var invalidId ClientError = ClientError{"Invalid ID"}

func onTrainDelete(db *gorm.DB, req *http.Request) (ResponseBody, HttpError) {
	queries := req.URL.Query()
	id := queries.Get("id")

	if id == "" {
		return nil, provideId
	}

	parsedId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, invalidId
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

	return nil, nil
}
