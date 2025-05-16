package main

import (
	"time"
)

func newSomeTime(t time.Time) optional[time.Time] {
	return optional[time.Time]{t, true}
}

func newNoneTime() optional[time.Time] {
	return optional[time.Time]{time.Time{}, false}
}

type expectedTimes struct {
	arrival, departure optional[time.Time]
}

func newExpectedTimes(arrival, departure optional[time.Time]) expectedTimes {
	return expectedTimes{arrival, departure}
}
