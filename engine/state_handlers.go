package main

func (s *EngineState) processRestart() {
	if s.Status == Restarting {
		s.responseOut <- NewEngineResponse(*s, NoOp)
		return
	}
	s.Status = Restarting
	s.responseOut <- NewEngineResponse(*s, Success)
	s.Status = Running
	s.responseOut <- NewEngineResponse(*s, Success)
}

func (s *EngineState) processPause() {
	if s.Status == Running {
		s.Status = Pausing
		s.responseOut <- NewEngineResponse(*s, Success)
		s.Status = Paused
		s.responseOut <- NewEngineResponse(*s, Success)
	} else {
		s.responseOut <- NewEngineResponse(*s, NoOp)
	}
}

func (s *EngineState) processUnpause() {
	if s.Status == Paused {
		s.Status = Unpausing
		s.responseOut <- NewEngineResponse(*s, Success)
		s.Status = Running
		s.responseOut <- NewEngineResponse(*s, Success)
	} else {
		s.responseOut <- NewEngineResponse(*s, NoOp)
	}
}

func (s *EngineState) processTrainEvent(events []TrainEvent) {
	for _, event := range events {
		switch event.EventType() {
		case CreateTrain:
			e := event.(EventCreateTrain)
			s.Trains[e.Name] = &e.Train
			s.responseOut <- NewEngineResponse(*s, Success)
		case DeleteTrain:
			e := event.(EventDeleteTrain)
			delete(s.Trains, e.name)
			s.responseOut <- NewEngineResponse(*s, Success)
		default:
			panic("unexpected engine.TrainEventType")
		}
	}
}

func (s *EngineState) processStationEvent(events []StationEvent) {
	for _, event := range events {
		switch event.EventType() {
		case CreateStation:
			e := event.(EventCreateStation)
			s.Stations[e.Name] = e.Station
			s.responseOut <- NewEngineResponse(*s, Success)
		case DeleteStation:
			e := event.(EventDeleteStation)
			delete(s.Trains, e.name)
			s.responseOut <- NewEngineResponse(*s, Success)
		default:
			panic("unexpected engine.StationEventType")
		}
	}
}
func (s *EngineState) processJourneyEvent(events []EventCreateJourney) {
	for _, event := range events {
		startStation, endStation := s.Stations[event.A], s.Stations[event.B]
		train := s.Trains[event.TrainName]
		s.Journeys = append(
			s.Journeys,
			newSimJourney(train, startStation, endStation),
		)
	}
	s.responseOut <- NewEngineResponse(*s, Success)
}
