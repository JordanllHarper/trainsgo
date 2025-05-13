package main

import (
	"fmt"
	"time"
)

type optionalTime struct {
	t      time.Time
	isSome bool
}

func NewSomeTime(t time.Time) optionalTime {
	return optionalTime{t, true}
}

func NewNoneTime() optionalTime {
	return optionalTime{time.Time{}, false}
}

type expectedTimes struct {
	arrival, departure optionalTime
}

func newExpectedTimes(arrival, departure optionalTime) expectedTimes {
	return expectedTimes{arrival, departure}
}

func (ot optionalTime) String() string {
	if ot.isSome {
		return fmt.Sprintf("%v", ot.t)
	} else {
		return fmt.Sprintf("None")
	}
}
