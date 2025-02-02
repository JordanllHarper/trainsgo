package main

import (
	"fmt"

	"github.com/JordanllHarper/trainsgo/backend/common"
	"github.com/JordanllHarper/trainsgo/backend/engine"
)

func main() {

	initialState := engine.NewEngineState([]common.Train{}, engine.Running)

	events := make(chan engine.Event)
	states := make(chan engine.EngineState)

	go engine.Run(events, states)

	for {
		go func() {
			events <- engine.NewEvent(initialState, engine.PauseSimulation)
		}()
		state := <-states
		fmt.Printf("%v\n", state.Status.PrettyPrint())
		state = <-states
		fmt.Printf("%v\n", state.Status.PrettyPrint())

	}

}
