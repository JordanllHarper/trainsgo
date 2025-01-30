package engine

import (
	"github.com/JordanllHarper/trainsgo/backend/common"
)

// The state we will send to consumers.
// All trains
type EngineState struct {
	Trains []common.Train
}
