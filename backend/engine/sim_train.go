package engine

import "github.com/JordanllHarper/trainsgo/backend/common"

type Journey struct {
	a common.Station
	b common.Station
}

type SimTrain struct {
	common.Train                   // internal state
	info_chan    chan common.Train // for external listeners
}
