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

	CreateStation StationEventType = iota
	DeleteStation
)

type (
	PlaybackEvent    int
	TrainEventType   int
	StationEventType int

	EventCreateTrain struct{ common.Train }
	EventDeleteTrain struct{ name string }

	EventCreateStation struct{ common.Station }
	EventDeleteStation struct{ name string }

	TrainEvent interface {
		EventType() TrainEventType
	}

	StationEvent interface {
		EventType() StationEventType
	}

	Event struct {
		pb      *PlaybackEvent
		train   TrainEvent
		station StationEvent
	}
)

// Create a new event to send to the simulation
// Provided so batch requests can be sent for more efficient transmission.
func NewEvent(pbEvents *PlaybackEvent, trainEvent TrainEvent, stationEvent StationEvent) Event {
	return Event{pbEvents, trainEvent, stationEvent}
}

//

// Train related functionality
func NewTrainEvent(e TrainEvent) Event { return Event{nil, e, nil} }

func NewEventCreateTrain(t common.Train) TrainEvent { return EventCreateTrain{t} }
func NewEventDeleteTrain(name string) TrainEvent    { return EventDeleteTrain{name} }

func (e EventCreateTrain) EventType() TrainEventType { return CreateTrain }
func (e EventDeleteTrain) EventType() TrainEventType { return DeleteTrain }

//

// Station related functionality
func NewStationEvent(e StationEvent) Event { return Event{nil, nil, e} }

func NewEventCreateStation(t common.Station) StationEvent { return EventCreateStation{t} }
func NewEventDeleteStation(name string) StationEvent      { return EventDeleteStation{name} }

func (e EventCreateStation) EventType() StationEventType { return CreateStation }
func (e EventDeleteStation) EventType() StationEventType { return DeleteStation }

//

// Playback control functionality
func NewPlaybackEvent(e PlaybackEvent) Event { return Event{pb: &e} }

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
