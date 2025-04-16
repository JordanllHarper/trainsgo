package engine

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

func (s *EngineState) processTrainEvent(event TrainEvent) {
	switch event.EventType() {

	case CreateTrain:
		e := event.(EventCreateTrain)
		s.Trains = append(s.Trains, e.Train)
		s.responseOut <- NewEngineResponse(*s, Success)
	case DeleteTrain:
		e := event.(EventDeleteTrain)
		for i, t := range s.Trains {
			if t.Name == e.name {
				s.Trains = common.RemoveIndexSlice(s.Trains, i)
				break
			}
		}
		s.responseOut <- NewEngineResponse(*s, Success)
	default:
		panic("unexpected engine.TrainEventType")
	}
}

func (s *EngineState) processStationEvent(event StationEvent) {
	switch event.EventType() {

	case CreateStation:
		e := event.(EventCreateStation)
		s.Stations = append(s.Stations, e.Station)
		s.responseOut <- NewEngineResponse(*s, Success)
	case DeleteStation:
		e := event.(EventDeleteStation)
		for i, t := range s.Stations {
			if t.Name == e.name {
				s.Stations = common.RemoveIndexSlice(s.Stations, i)
				break
			}
		}
		s.responseOut <- NewEngineResponse(*s, Success)
	default:
		panic("unexpected engine.StationEventType")
	}
}
