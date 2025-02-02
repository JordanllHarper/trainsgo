package engine

import "fmt"

/*
TODO:
Setup registering new stations.
Setup registering new trains
Simulate all our trains moving - incrementing and decrementing coordinates
*/

type PlaybackEvents int

const (
	// Will refresh the simulation and move back to a starting state or restart
	RestartSimulation PlaybackEvents = iota
	PauseSimulation
	UnpauseSimulation
	QuitSimulation
)

type Event struct {
	currentEngineState EngineState
	pbEvents           PlaybackEvents
}

func NewEvent(currentEngineState EngineState, pbEvents PlaybackEvents) Event {
	return Event{
		currentEngineState, pbEvents,
	}
}

func Run(in_events chan Event, state_out chan EngineState) error {
	for {
		x := <-in_events
		// fmt.Printf("Received event %v", x)

		switch s := x.currentEngineState; x.pbEvents {
		case PauseSimulation:
			s.processPause(state_out)
		case QuitSimulation:
			fmt.Println("Quitting...")
			return nil
		case RestartSimulation:
			s.processRestart(state_out)
		case UnpauseSimulation:
			s.processUnpause(state_out)
		default:
			panic(fmt.Sprintf("Unexpected event: %#v", x.pbEvents))
		}

	}

}
