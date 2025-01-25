package api

import "net/http"

type TrainResponseSingular struct {
	TrainEntity
	statusCode int
}

func NewTrainResponseSingular(statusCode int, train TrainEntity) TrainResponseSingular {
	return TrainResponseSingular{TrainEntity: train, statusCode: http.StatusOK}
}

func (response TrainResponseSingular) StatusCode() int { return response.statusCode }

type TrainResponseMultiple struct {
	Trains     []TrainEntity
	statusCode int
}

func NewTrainResponseMultiple(statusCode int, trains []TrainEntity) TrainResponseMultiple {
	return TrainResponseMultiple{Trains: trains, statusCode: http.StatusOK}
}

func (response TrainResponseMultiple) StatusCode() int { return response.statusCode }
