package main

import (
	"time"
)

type expectedTimes struct {
	Arrival   *time.Time `json:"arrival"`
	Departure *time.Time `json:"departure"`
}

func newExpectedTimes(arrival, departure *time.Time) expectedTimes {
	return expectedTimes{arrival, departure}
}
