package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/JordanllHarper/trainsgo/backend/engine"
)

func monitorAndPrintState(states chan engine.EngineState) {
	for {
		state := <-states
		fmt.Printf("STATE: %v\n", state.Status.PrettyPrint())
	}
}

func main() {

	events := make(chan engine.Event)
	states := make(chan engine.EngineState)

	fmt.Println("Press:\n 1 -> pause\n 2 -> restart\n q -> quit")

	go engine.Run(events, states)
	go monitorAndPrintState(states)

	reader := bufio.NewReader(os.Stdin)
	run := true
	for run {
		text, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalln("Error:", err)
		}

		switch strings.TrimSpace(text) {
		case "1":
			events <- engine.NewPlaybackEvent(engine.PauseSimulation)
		case "2":
			events <- engine.NewPlaybackEvent(engine.RestartSimulation)
		case "q":
			run = false
		}
	}
	print("Quitting...")
}
