package main

type StationEntity struct {
	DbFields
	Station
}

type Station struct {
	Name string `json:"name"`
	Coordinates
}
