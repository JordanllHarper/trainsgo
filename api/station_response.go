package main

type StationResponseSingular struct {
	StationEntity
	statusCode int
}

func NewStationResponseSingular(statusCode int, station StationEntity) StationResponseSingular {
	return StationResponseSingular{station, statusCode}
}

func (response StationResponseSingular) StatusCode() int { return response.statusCode }

type StationResponseMultiple struct {
	Stations   []StationEntity
	statusCode int
}

func NewStationResponseMultiple(statusCode int, station []StationEntity) StationResponseMultiple {
	return StationResponseMultiple{station, statusCode}
}

func (response StationResponseMultiple) StatusCode() int { return response.statusCode }
