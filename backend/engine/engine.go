package engine

import (
	"fmt"

	"github.com/JordanllHarper/trainsgo/backend/common"
)

/*
TODO:
Setup registering new stations.
Setup registering new trains
Simulate all our trains moving - incrementing and decrementing coordinates
*/

func log(message string) {
	fmt.Println("ENGINE:", message)
}

func handlePlaybackEvent(pbEvent PlaybackEvent, currentState *EngineState, stateOut chan EngineState) bool {

	switch pbEvent {
	case PauseSimulation:
		currentState.processPause(stateOut)
	case QuitSimulation:
		log("Quitting simulation...")
		return false
	case RestartSimulation:
		currentState.processRestart(stateOut)
	case UnpauseSimulation:
		currentState.processUnpause(stateOut)
	default:
		panic(fmt.Sprintf("Unexpected playback event: %#v", pbEvent.pretty()))
	}
	return true
}

func Run(inEvents chan Event, stateOut chan EngineState) error {

	currentState := NewEngineState([]common.Train{}, Running)
	stateOut <- currentState
	run := true
	for run {
		event := <-inEvents
		if event.PlaybackEvent != nil {
			pbEvent := *event.PlaybackEvent
			run = handlePlaybackEvent(pbEvent, &currentState, stateOut)
		}
		if event.TrainEvent != nil {
			tEvent := event.TrainEvent

		}
	}
	return nil
}
