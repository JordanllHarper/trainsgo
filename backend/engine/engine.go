package engine

/*
TODO:
Setup registering new stations.
Setup registering new trains
Simulate all our trains moving - incrementing and decrementing coordinates
*/

type EventType int

const (
	// Will refresh the simulation and move back to a starting state or restart
	RestartSimulation EventType = iota
	PauseSimulation
	UnpauseSimulation
)

type Event struct {
	EventType
}

func Run(events chan Event) chan EngineState {
	engineChan := make(chan EngineState)

	return engineChan
}
