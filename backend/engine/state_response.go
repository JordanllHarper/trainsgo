package engine

import "fmt"

type ResponseCode int

const (
	Success ResponseCode = iota
	Initialised
	NoOp            // for when the command will do nothing
	InvalidCreation // when the user tries to create an invalid entity
	InvalidDeletion // when the user tries to delete an entity that they shouldn't be able to (e.g. a station that has a train travelling to it)
)

type EngineResponse struct {
	EngineState
	ResponseCode
}

func NewEngineResponse(state EngineState, r ResponseCode) EngineResponse {
	return EngineResponse{state, r}
}

func (rc ResponseCode) ToString() string {
	switch rc {
	case Initialised:
		return "Initialised"
	case InvalidCreation:
		return "Invalid Creation"
	case InvalidDeletion:
		return "Invalid Deletion"
	case NoOp:
		return "No Operation"
	case Success:
		return "Success"
	default:
		panic(fmt.Sprintf("unexpected ResponseCode: %#v", rc))
	}
}
