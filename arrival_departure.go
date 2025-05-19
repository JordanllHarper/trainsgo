package main

import (
	"time"
)

type ExpectedTimes struct {
	Departure time.Time `json:"departure"`
	Arrival   time.Time `json:"arrival"`
}

func newExpectedTimes(arrival, departure time.Time) ExpectedTimes {
	return ExpectedTimes{
		Arrival:   arrival,
		Departure: departure,
	}
}
