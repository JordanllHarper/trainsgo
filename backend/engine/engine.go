package engine

import (
	"fmt"
)

/*
	TODO:
	Setup registering new stations.
	Simulate all our trains moving - incrementing and decrementing coordinates
*/

func log(message string) {
	fmt.Println("ENGINE:", message)
}

func handlePlaybackEvent(pbEvent PlaybackEvent, currentState *EngineState) bool {

	switch pbEvent {
	case PauseSimulation:
		currentState.processPause()
	case QuitSimulation:
		log("Quitting simulation...")
		return false
	case RestartSimulation:
		currentState.processRestart()
	case UnpauseSimulation:
		currentState.processUnpause()
	default:
		panic(fmt.Sprintf("Unexpected playback event: %#v", pbEvent.Pretty()))
	}
	return true
}

func Run(inEvents chan Event, stateOut chan EngineResponse) error {
	currentState := NewEngineState(Running, stateOut)
	stateOut <- NewEngineResponse(currentState, Initialised)
	run := true
	for run {
		event := <-inEvents
		if event.pb != nil {
			pbEvent := *event.pb
			run = handlePlaybackEvent(pbEvent, &currentState)
		}
		if event.train != nil {
			currentState.processTrainEvent(event.train)
		}
		if event.station != nil {
			currentState.processStationEvent(event.station)
		}
	}
	return nil
}
