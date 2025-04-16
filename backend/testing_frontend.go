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

func monitorAndPrintState(responses chan engine.EngineResponse) {
	for {
		response := <-responses
		fmt.Printf("STATE: %v\n", response.Status.ToString())
		fmt.Printf("TRAINS: %v\n", common.ArrayToString(response.Trains))
		fmt.Printf("STATIONS: %v\n", common.ArrayToString(response.Stations))
		fmt.Printf("CODE: %v\n", response.ResponseCode.ToString())
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
	responses := make(chan engine.EngineResponse)

	help := generateHelp([]help{
		{"p", "[P]lay/Pause"},
		{"r", "[R]estart"},
		{"ct", "[C]reate new test [t]rain"},
		{"cs", "[C]reate new test [s]tation*s*"},
		// TODO: {"cj", "[C]reate new test [j]ourney"},
		{"d", "[D]elete test train"},
		{"q", "[Q]uit"},
		{"h", "[H]elp"},
	},
	)

	go engine.Run(events, responses)
	go monitorAndPrintState(responses)

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
		case "ct":
			events <- engine.NewTrainEvent(engine.NewEventCreateTrain(
				common.NewTrain("test", 3, common.Coordinates{X: 0, Y: 0}, common.Unused)),
			)

		case "cs":
			events <- engine.NewStationEvent(engine.NewEventCreateStation(
				common.NewStation("test", common.Coordinates{X: 0, Y: 0})),
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
