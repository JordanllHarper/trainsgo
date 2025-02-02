package main

import "github.com/JordanllHarper/trainsgo/backend/common"

type StationEntity struct {
	DbFields
	common.Station
}

type TrainEntity struct {
	DbFields
	common.Train
}
