package main

import "github.com/JordanllHarper/trainsgo/backend/common"

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
