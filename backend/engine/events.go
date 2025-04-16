package engine

import (
	"fmt"

	"github.com/JordanllHarper/trainsgo/backend/common"
)

const (
	// Will refresh the simulation and move back to a starting state or restart
	RestartSimulation PlaybackEvent = iota
	PauseSimulation
	UnpauseSimulation
	QuitSimulation

	CreateTrain TrainEventType = iota
	DeleteTrain
)

type (
	PlaybackEvent  int
	TrainEventType int

	EventCreateTrain struct{ common.Train }
	EventDeleteTrain struct{ name string }

	TrainEvent interface {
		EventType() TrainEventType
		ToString() string
	}

	Event struct {
		*PlaybackEvent
		TrainEvent
	}
)

func (e EventCreateTrain) EventType() TrainEventType { return CreateTrain }
func (e EventDeleteTrain) EventType() TrainEventType { return DeleteTrain }

func NewEventCreateTrain(t common.Train) TrainEvent { return EventCreateTrain{t} }
func NewEventDeleteTrain(name string) TrainEvent    { return EventDeleteTrain{name} }

func (event PlaybackEvent) Pretty() string {
	var s string
	switch event {
	case PauseSimulation:
		s = "Event - Pause Simulation"
	case QuitSimulation:
		s = "Event - Quit Simulation"
	case RestartSimulation:
		s = "Event - Restart Simulation"
	case UnpauseSimulation:
		s = "Event - Unpause Simulation"
	default:
		panic(fmt.Sprintf("unexpected engine.PlaybackEvents num: %#v", event))
	}
	return s
}

func NewEvent(pbEvents *PlaybackEvent, trainEvent TrainEvent) Event {
	return Event{pbEvents, trainEvent}
}

func NewPlaybackEvent(e PlaybackEvent) Event { return Event{&e, nil} }
func NewTrainEvent(e TrainEvent) Event       { return Event{nil, e} }

func (e EventCreateTrain) ToString() string { return fmt.Sprintf("%v", e) }
func (e EventDeleteTrain) ToString() string { return fmt.Sprintf("%v", e) }
