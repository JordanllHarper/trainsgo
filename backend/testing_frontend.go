package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/JordanllHarper/trainsgo/backend/common"
	"github.com/JordanllHarper/trainsgo/backend/engine"
)

func monitorAndPrintState(states chan engine.EngineState) {
	for {
		state := <-states
		fmt.Printf("STATE: %v\n", state.Status.ToString())
		fmt.Printf("TRAINS: %v\n", common.ArrayToString(state.Trains))
	}
}

type help struct {
	key, command string
}

func generateHelp(opts []help) string {
	s := "Press:"
	for _, v := range opts {
		s = fmt.Sprintf("%s\n%s -> %s", s, v.key, v.command)
	}
	return s
}

func main() {

	events := make(chan engine.Event)
	states := make(chan engine.EngineState)

	help := generateHelp([]help{
		{"p", "[P]lay/Pause"},
		{"r", "[R]estart"},
		{"c", "[C]reate new test train"},
		{"d", "[D]elete test train"},
		{"q", "[Q]uit"},
		{"h", "[H]elp"},
	},
	)

	go engine.Run(events, states)
	go monitorAndPrintState(states)

	reader := bufio.NewReader(os.Stdin)
	run := true

	fmt.Println(help)
	for run {
		text, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalln("Error:", err)
		}

		switch strings.TrimSpace(text) {
		case "p":
			events <- engine.NewPlaybackEvent(engine.PauseSimulation)
		case "r":
			events <- engine.NewPlaybackEvent(engine.RestartSimulation)
		case "c":
			events <- engine.NewTrainEvent(engine.NewEventCreateTrain(
				common.NewTrain("test", 3, common.Coordinates{X: 0, Y: 0}, common.Unused)),
			)
		case "d":
			events <- engine.NewTrainEvent(engine.NewEventDeleteTrain("test"))
		case "q":
			run = false
		case "h":
			fmt.Println(help)
		}
	}
	fmt.Println("Quitting...")
}
