package main

import "time"

type expectedTimes struct {
	arrival, departure time.Time
}

func newExpectedTimes(arrival, departure time.Time) expectedTimes {
	return expectedTimes{arrival, departure}
}
