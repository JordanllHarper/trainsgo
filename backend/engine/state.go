package engine

import (
	"fmt"

	"github.com/JordanllHarper/trainsgo/backend/common"
)

const (
	Restarting EngineStatus = iota
	Running
	Pausing
	Paused
	Unpausing
)

type (
	EngineStatus int

	// The state we will send to consumers.
	EngineState struct {
		Trains      []common.Train
		Stations    []common.Station
		Status      EngineStatus
		responseOut chan EngineResponse
	}
)

func NewEngineState(status EngineStatus, stateOut chan EngineResponse) EngineState {
	return EngineState{[]common.Train{}, []common.Station{}, status, stateOut}
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
	default:
		panic(fmt.Sprintf("unexpected engine.EngineStatus num: %#v", s))
	}
}
