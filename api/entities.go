package main

import "github.com/JordanllHarper/trainsgo/common"

type (
	StationEntity struct {
		DbFields
		common.Station
	}

	TrainEntity struct {
		DbFields
		common.Train
	}
)
