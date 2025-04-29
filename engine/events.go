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
	PlaybackEvent int

	TrainEventType   int
	StationEventType int
	JourneyEventType int

	EventCreateTrain struct{ common.Train }
	EventDeleteTrain struct{ name string }

	EventCreateStation struct{ common.Station }
	EventDeleteStation struct{ name string }

	EventCreateJourney struct{ common.Journey }

	TrainEvent interface {
		EventType() TrainEventType
	}

	StationEvent interface {
		EventType() StationEventType
	}

	Event struct {
		pb      *PlaybackEvent
		train   []TrainEvent
		station []StationEvent
		journey []EventCreateJourney
	}
)

// Create a new event to send to the simulation
// Provided so batch requests can be sent for more efficient transmission.
func NewEvent(pbEvent *PlaybackEvent, tEvent []TrainEvent, sEvent []StationEvent, jEvent []EventCreateJourney) Event {
	return Event{pbEvent, tEvent, sEvent, jEvent}
}

//

// Train related functionality
func NewTrainEvent(e TrainEvent) Event    { return Event{train: []TrainEvent{e}} }
func NewTrainEvents(e []TrainEvent) Event { return Event{train: e} }

func NewEventCreateTrain(t common.Train) TrainEvent { return EventCreateTrain{t} }
func NewEventDeleteTrain(name string) TrainEvent    { return EventDeleteTrain{name} }

func (e EventCreateTrain) EventType() TrainEventType { return CreateTrain }
func (e EventDeleteTrain) EventType() TrainEventType { return DeleteTrain }

//

// Station related functionality
func NewStationEvent(e StationEvent) Event    { return Event{station: []StationEvent{e}} }
func NewStationEvents(e []StationEvent) Event { return Event{station: e} }

func NewEventCreateStation(s common.Station) StationEvent { return EventCreateStation{s} }
func NewEventDeleteStation(name string) StationEvent      { return EventDeleteStation{name} }

func (e EventCreateStation) EventType() StationEventType { return CreateStation }
func (e EventDeleteStation) EventType() StationEventType { return DeleteStation }

// Journey functionality

func NewJourneyEvent(e EventCreateJourney) Event    { return Event{journey: []EventCreateJourney{e}} }
func NewJourneyEvents(e []EventCreateJourney) Event { return Event{journey: e} }

func NewEventCreateJourney(j common.Journey) EventCreateJourney { return EventCreateJourney{j} }

func (e EventCreateJourney) EventType() StationEventType { return CreateStation }

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
