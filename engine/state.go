package main

import (
	"fmt"

	"github.com/JordanllHarper/trainsgo/common"
)

const (
	Restarting EngineStatus = iota
	Initialised
	Running
	Pausing
	Paused
	Unpausing
)

type (
	EngineStatus int

	// The state we will send to consumers.
	EngineState struct {
		Trains   map[string]*common.Train
		Stations map[string]common.Station
		Journeys []simJourney
		Status   EngineStatus
		// internal so listeners can't publish responses unless they have a reference to their channel
		responseOut chan EngineResponse
		eventOut    chan EngineEvent
	}
)

func NewEngineState(status EngineStatus, stateOut chan EngineResponse, eventOut chan EngineEvent) EngineState {
	return EngineState{map[string]*common.Train{}, map[string]common.Station{}, []simJourney{}, status, stateOut, eventOut}
}

func (s EngineStatus) ToString() string {
	switch s {
	case Paused:
		return "Paused"
	case Pausing:
		return "Pausing"
	case Restarting:
		return "Restarting"
	case Running:
		return "Running"
	case Unpausing:
		return "Unpausing"
	case Initialised:
		return "Initialised"
	}
	panic(fmt.Sprintf("unexpected engine.EngineStatus num: %#v", s))
}
