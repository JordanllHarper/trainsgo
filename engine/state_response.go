package main

import "fmt"

type ResponseCode int
type EventCode int

const (
	Success         ResponseCode = iota
	NoOp                         // for when the command will do nothing
	InvalidCreation              // when the user tries to create an invalid entity
	InvalidDeletion              // when the user tries to delete an entity that they shouldn't be able to (e.g. a station that has a train travelling to it)

	Started EventCode = iota
	Updated
)

type EngineResponse struct {
	EngineState
	ResponseCode
}

type EngineEvent struct {
	EngineState
	EventCode
}

func NewEngineEvent(state EngineState, e EventCode) EngineEvent {
	return EngineEvent{state, e}
}

func NewEngineResponse(state EngineState, r ResponseCode) EngineResponse {
	return EngineResponse{state, r}
}

func (r ResponseCode) ToString() string {
	switch r {
	case InvalidCreation:
		return "Invalid Creation"
	case InvalidDeletion:
		return "Invalid Deletion"
	case NoOp:
		return "No Operation"
	case Success:
		return "Success"
	default:
		panic(fmt.Sprintf("unexpected ResponseCode: %#v", r))
	}
}

func (e EventCode) ToString() string {
	switch e {
	case Started:
		return "Started"
	case Updated:
		return "Updated"
	default:
		panic(fmt.Sprintf("unexpected EventCode: %#v", e))
	}
}
