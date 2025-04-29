package main

import "github.com/JordanllHarper/trainsgo/common"

type SimTrain struct {
	common.Train                   // internal state
	info_chan    chan common.Train // for external listeners
}

type simJourney struct {
	*common.Train
	A, B common.Station
}

func newSimJourney(t *common.Train, a, b common.Station) simJourney {
	return simJourney{t, a, b}
}
