package main

import (
	"time"
)

type expectedTimes struct {
	Departure time.Time `json:"departure"`
	Arrival   time.Time `json:"arrival"`
}

func newExpectedTimes(arrival, departure time.Time) expectedTimes {
	return expectedTimes{
		Arrival:   arrival,
		Departure: departure,
	}
}
