package engine

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/JordanllHarper/trainsgo/backend/common"
)

/*
	TODO:
	Setup registering new stations.
	Simulate all our trains moving - incrementing and decrementing coordinates
*/

func handlePlaybackEvent(pbEvent PlaybackEvent, currentState *EngineState) bool {

	switch pbEvent {
	case PauseSimulation:
		currentState.processPause()
	case QuitSimulation:
		log.Default().Println("Quitting simulation...")
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

func computeAxis(current, dest, speed int) int {
	switch {
	case dest > current:
		difference := dest - current
		if speed > difference {
			return dest
		}

		return current + speed
	case dest < current:
		difference := current - dest
		if speed > difference {
			return dest
		}
		return current - speed
	default:
		return dest
	}
}

func moveTrain(t *common.Train, to common.Coordinates) {
	speed := t.TopSpeed
	newX, newY := computeAxis(t.X, to.X, speed), computeAxis(t.Y, to.Y, speed)
	t.Coordinates = common.Coordinates{X: newX, Y: newY}
	if t.Coordinates == to {
		t.Status = common.Arrived
	}
}

func update(currentState *EngineState, ctx context.Context) error {
	currentState.Status = Running
	for {
		time.Sleep(1 * time.Second)
		select {
		case <-ctx.Done():
			return nil
		default:
		}
		finishedJourneyIndexes := []int{}
		for i, j := range currentState.Journeys {
			train, from, to := j.Train, j.A, j.B
			switch train.Status {
			case common.Unused:
				train.Coordinates = from.Coordinates
				train.Status = common.Travelling
				fallthrough
			case common.Travelling:
				moveTrain(train, to.Coordinates)
				currentState.eventOut <- NewEngineEvent(*currentState, Updated)
				if train.Status != common.Arrived {
					continue
				}
			case common.Arrived:
				finishedJourneyIndexes = append(finishedJourneyIndexes, i) // we will bulk clean this up
				currentState.eventOut <- NewEngineEvent(*currentState, Updated)
				train.Status = common.Unused
				currentState.eventOut <- NewEngineEvent(*currentState, Updated)
			default:
				panic(fmt.Sprintf("unexpected common.Status: %#v", train.Status))
			}
		}
		for _, j := range finishedJourneyIndexes {
			newJourneys := common.RemoveIndexSliceCopied(currentState.Journeys, j)
			currentState.Journeys = newJourneys
		}
	}
}

func Run(inEvents chan Event, outResponses chan EngineResponse, outEvents chan EngineEvent) error {
	currentState := NewEngineState(Initialised, outResponses, outEvents)
	outResponses <- NewEngineResponse(currentState, Success)
	ctx, cancel := context.WithCancel(context.Background())
	go update(&currentState, ctx)
	for {
		event := <-inEvents
		if event.pb != nil {
			pbEvent := *event.pb
			shouldContinue := handlePlaybackEvent(pbEvent, &currentState)
			if !shouldContinue {
				cancel()
				break
			}
		}
		if event.train != nil {
			currentState.processTrainEvent(event.train)
		}
		if event.station != nil {
			currentState.processStationEvent(event.station)
		}
		if event.journey != nil {
			currentState.processJourneyEvent(event.journey)
		}

	}
	return nil
}
