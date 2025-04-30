package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/JordanllHarper/trainsgo/common"
)

func printStateResponse(r EngineResponse) {
	fmt.Printf("STATE: %v\n", r.Status.ToString())
	fmt.Printf("TRAINS: %v\n", common.MapToString(r.Trains))
	fmt.Printf("STATIONS: %v\n", common.MapToString(r.Stations))
	fmt.Printf("JOURNEYS: %v\n", common.SliceToString(r.Journeys))
	fmt.Printf("CODE: %v\n", r.ResponseCode.ToString())
}

func printStateEvent(e EngineEvent) {
	fmt.Println("!!! EVENT !!!")
	fmt.Printf("STATE: %v\n", e.Status.ToString())
	fmt.Printf("TRAINS: %v\n", common.MapToString(e.Trains))
	fmt.Printf("STATIONS: %v\n", common.MapToString(e.Stations))
	fmt.Printf("JOURNEYS: %v\n", common.SliceToString(e.Journeys))
	fmt.Printf("CODE: %v\n", e.EventCode.ToString())
}

func monitorAndPrintResponse(responses chan EngineResponse) {
	for {
		r := <-responses
		printStateResponse(r)
	}
}

func monitorAndPrintEvent(events chan EngineEvent) {
	for {
		e := <-events
		printStateEvent(e)
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

	inEvents := make(chan Event)
	outResponses := make(chan EngineResponse)
	outEvents := make(chan EngineEvent)

	go run(inEvents, outResponses, outEvents)
	go monitorAndPrintResponse(outResponses)
	go monitorAndPrintEvent(outEvents)

	select {
	case r := <-outResponses:
		if r.ResponseCode == Success {
			printStateResponse(r)
			break
		} else {
			log.Fatalf("Response code was not success, exiting %v\n", r.ResponseCode.ToString())
		}
	}

	err := loop(inEvents)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Println("Quitting...")
}

var h string = generateHelp([]help{
	{"p", "[P]layback"}, // pl for unpause, pa for pause, r for restart
	{"ct", "[C]reate new test [t]rain"},
	{"cs", "[C]reate new test [s]tation*s*"},
	{"cj", "[C]reate new test [j]ourney"},
	{"d", "[D]elete test train"},
	{"q", "[Q]uit"},
	{"h", "[H]elp"},
},
)

func handlePlaybackSelection(args []string) (Event, error) {
	subcommand, err := safeSliceAccess(args, 1)
	if err != nil {
		return Event{}, errors.New("Invalid playback subcommand")
	}

	var event Event
	switch subcommand {
	case "pl":
		event = NewPlaybackEvent(UnpauseSimulation)
	case "pa":
		event = NewPlaybackEvent(PauseSimulation)
	case "r":
		event = NewPlaybackEvent(RestartSimulation)

	}

	return event, nil
}

func loop(events chan Event) error {
	fmt.Println("Press h for [H]elp")
	for {
		text, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			log.Fatalln("Error:", err)
		}

		commandArgs := strings.Split(strings.TrimSpace(text), " ")
		command, err := safeSliceAccess(commandArgs, 0)
		if err != nil {
			log.Println("Did not provide a command, please try again")
			continue
		}

		switch command {
		case "p":
			pb, err := handlePlaybackSelection(commandArgs[1:])
			if err != nil {
				log.Println("Did not provide a suitable subcommand, please try again")
				continue
			}
			events <- pb
		case "ct":
			events <- NewTrainEvent(NewEventCreateTrain(
				common.NewTrain("testTrain", 1, common.Coordinates{X: 0, Y: 0}, common.Unused)),
			)
		case "cs":
			events <- NewStationEvents(
				[]StationEvent{
					NewEventCreateStation(common.NewStation("testStation", common.Coordinates{X: 0, Y: 0})),
					NewEventCreateStation(common.NewStation("testStation2", common.Coordinates{X: 10, Y: 10})),
				},
			)
		case "cj":
			events <- NewJourneyEvent(
				NewEventCreateJourney(
					common.NewJourney("testStation", "testStation2", "testTrain"),
				),
			)
		case "d":
			events <- NewTrainEvent(NewEventDeleteTrain("testTrain"))
		case "q":
			return nil
		case "h":
			fmt.Println(h)
		default:
			fmt.Println("Unrecognised command, please try again or use [h]elp")
		}
	}
}
