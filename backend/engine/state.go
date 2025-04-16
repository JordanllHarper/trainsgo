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
	Trains   []common.Train
	Status   EngineStatus
	stateOut chan EngineState
}

func NewEngineState(trains []common.Train, status EngineStatus, stateOut chan EngineState) EngineState {
	return EngineState{trains, status, stateOut}
}

func (s *EngineState) processRestart() {
	s.Status = Restarting
	s.stateOut <- *s
	s.Status = Running
	s.stateOut <- *s
}

func (s *EngineState) processPause() {
	s.Status = Pausing
	s.stateOut <- *s
	s.Status = Paused
	s.stateOut <- *s
}

func (s *EngineState) processUnpause() {
	s.Status = Unpausing
	s.stateOut <- *s
	s.Status = Running
	s.stateOut <- *s
}

func (s *EngineState) processTrainEvent(event TrainEvent) {
	switch event.EventType() {

	case CreateTrain:
		e := event.(EventCreateTrain)
		s.Trains = append(s.Trains, e.Train)
		s.stateOut <- *s
	case DeleteTrain:
		e := event.(EventDeleteTrain)
		for i, t := range s.Trains {
			if t.Name == e.name {
				s.Trains = common.RemoveIndexSlice(s.Trains, i)
				break
			}
		}
		s.stateOut <- *s
	default:
		panic("unexpected engine.TrainEventType")
	}
}
