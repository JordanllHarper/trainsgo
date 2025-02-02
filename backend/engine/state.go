package engine

import (
	"github.com/JordanllHarper/trainsgo/backend/common"
)

type EngineStatus int

const (
	Restarting EngineStatus = iota
	Running
	Pausing
	Paused
	Unpausing
	Unpaused
)

func (s EngineStatus) PrettyPrint() string {
	switch s {
	case Paused:
		return "Paused"
	case Pausing:
		return "Pausing"
	case Restarting:
		return "Restarting"
	case Running:
		return "Running"
	case Unpaused:
		return "Unpaused"
	case Unpausing:
		return "Unpausing"
	}

	return "Not an engine status :P"
}

// The state we will send to consumers.
// All trains
type EngineState struct {
	Trains []common.Train
	Status EngineStatus
}

func NewEngineState(trains []common.Train, status EngineStatus) EngineState {
	return EngineState{trains, status}
}

func (e EngineState) processRestart(stateOut chan EngineState) {
	stateOut <- NewEngineState([]common.Train{}, Restarting)
	// this might get more complicated if we need to restart trains, but we'll just clear them for now
	stateOut <- NewEngineState([]common.Train{}, Running)
}

func (e EngineState) processPause(stateOut chan EngineState) {
	stateOut <- NewEngineState(e.Trains, Pausing)
	// Some logic
	stateOut <- NewEngineState(e.Trains, Paused)
}

func (e EngineState) processUnpause(stateOut chan EngineState) {
	stateOut <- NewEngineState(e.Trains, Unpausing)
	// Some logic
	stateOut <- NewEngineState(e.Trains, Unpaused)
}
