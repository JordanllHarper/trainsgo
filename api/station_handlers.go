package main

import (
	"encoding/json"
	"fmt"
	"github.com/JordanllHarper/trainsgo/backend/common"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

func onStationGet(db *gorm.DB, req *http.Request) (ResponseBody, HttpError) {
	queries := req.URL.Query()
	id := queries.Get("id")

	if stringEmpty(id) {
		var stationEntities []StationEntity
		db.Find(&stationEntities)
		return NewStationResponseMultiple(http.StatusOK, stationEntities), nil
	}

	parsedId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, provideId()
	}

	var stationEntity StationEntity
	result := db.Where(&StationEntity{DbFields: DbFields{ID: uint(parsedId)}}).Find(&stationEntity)
	if result.RowsAffected == 0 {
		return NewEmptyResponseBody(), nil
	}

	return NewStationResponseSingular(http.StatusOK, stationEntity), nil
}

func onStationPost(db *gorm.DB, req *http.Request) (ResponseBody, HttpError) {
	body := req.Body
	if body == nil {
		return nil, needBody()
	}

	var station common.Station
	if err := json.NewDecoder(body).Decode(&station); err != nil {
		return nil, malformedBody()
	}

	sEntity := StationEntity{DbFields: DbFields{}, Station: station}
	db = db.Create(&sEntity)
	fmt.Printf("[INFO]: Inserted train entity: %v\n", sEntity)

	return NewStationResponseSingular(http.StatusCreated, sEntity), nil
}

func onStationPut(db *gorm.DB, req *http.Request) (ResponseBody, HttpError) {
	body := req.Body
	queries := req.URL.Query()
	id := queries.Get("id")
	if stringEmpty(id) {
		return nil, provideId()
	}

	parsedId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, invalidId()
	}

	var updatedStationFields common.Station
	if err := json.NewDecoder(body).Decode(&updatedStationFields); err != nil {
		return nil, malformedBody()
	}

	stationEntity := StationEntity{DbFields: DbFields{ID: uint(parsedId)}, Station: updatedStationFields}
	result := db.Updates(&stationEntity)
	if result.RowsAffected != 0 {
		fmt.Printf("[INFO]: Modified station entity id: %v\n", parsedId)
		return NewStationResponseSingular(http.StatusOK, stationEntity), nil
	} else {
		fmt.Printf("[INFO]: Attempted modification with no result: %v\n", parsedId)
		return NewEmptyResponseBody(), nil
	}
}

func onStationDelete(db *gorm.DB, req *http.Request) (ResponseBody, HttpError) {
	queries := req.URL.Query()
	id := queries.Get("id")
	if stringEmpty(id) {
		return nil, provideId()
	}

	parsedId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, invalidId()
	}

	trainEntity := &StationEntity{
		DbFields: DbFields{ID: uint(parsedId)},
	}
	result := db.Delete(&trainEntity)
	if result.RowsAffected != 0 {
		fmt.Printf("[INFO]: Deleted train entity id: %v\n", parsedId)
	} else {
		fmt.Printf("[INFO]: Attempted deletion with no result: %v\n", parsedId)
	}

	return NewEmptyResponseBody(), nil
}
