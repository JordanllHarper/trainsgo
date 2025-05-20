package main

import "net/http"

type (
	HttpResponse interface {
		HttpCode() int
		Body() any
	}

	statusOK      struct{ body any }
	statusCreated struct{ body any }

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
		Name       string `json:"name"`
		StationOne string `json:"stationOne"`
		StationTwo string `json:"stationTwo"`
	}

	tripPostBody struct {
		FromStationId  string        `json:"fromStationId"`
		ToStationId    string        `json:"toStationId"`
		TrainId        string        `json:"trainId"`
		StartingStatus TripStatus    `json:"startingStatus"`
		ExpTimes       ExpectedTimes `json:"expectedTimes"`
	}

	tripPutBody struct {
		Id        string     `json:"id"`
		NewStatus TripStatus `json:"status"`
	}

	renameBody struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	}

	deleteBody struct {
		Id string `json:"id"`
	}
)

func (s statusOK) HttpCode() int { return http.StatusOK }
func (s statusOK) Body() any     { return nil }

func (s statusCreated) HttpCode() int { return http.StatusCreated }
func (s statusCreated) Body() any     { return s.body }
