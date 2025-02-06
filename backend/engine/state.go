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

func (s *EngineState) processRestart(stateOut chan EngineState) {
	s.Status = Restarting
	stateOut <- *s
	s.Status = Running
	stateOut <- *s
}

func (s *EngineState) processPause(stateOut chan EngineState) {
	s.Status = Pausing
	stateOut <- *s
	s.Status = Paused
	stateOut <- *s
}

func (s *EngineState) processUnpause(stateOut chan EngineState) {
	s.Status = Unpausing
	stateOut <- *s
	s.Status = Running
	stateOut <- *s
}

func (s *EngineState) processTrainEvent(event TrainEvent) {
	switch event.EventType() {

	case CreateTrain:
		e := event.(EventCreateTrain)

		s.Trains = append(s.Trains, e.TrainToAdd)
		// handle creating a train
	case DeleteTrain:
		// handle deleting a train
	default:
		panic("unexpected engine.TrainEventType")
	}
}
