package main

import "fmt"

type (
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

	errorBody struct {
		Message string `json:"message"`
	}
)

func errMalformedBody() errorBody {
	return errorBody{Message: "Malformed body"}
}

func errBadId(id string) errorBody {
	return errorBody{Message: fmt.Sprintf("Bad ID format for %s", id)}
}

func errIdDoesntExist(id Id) errorBody {
	return errorBody{Message: fmt.Sprintf(
		"ID %s doesn't exist",
		id.String(),
	)}
}

func errIdExists(id Id) errorBody {
	return errorBody{
		Message: fmt.Sprintf(
			"%s ID already exists",
			id.String(),
		),
	}
}
