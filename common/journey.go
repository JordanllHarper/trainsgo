package common

type Journey struct {
	A, B, TrainName string
}

func NewJourney(a, b, trainName string) Journey {
	return Journey{a, b, trainName}
}
