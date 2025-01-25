package api

type Status string

const (
	Travelling   Status = "Travelling"
	Transferring        = "Transferring"
	Unused              = "Unused"
	Emergency           = "Emergency"
)

type TrainEntity struct {
	DbFields
	Train
}
type Train struct {
	Name string `json:"name"`
	// meters / second
	TopSpeed int `json:"top_speed"`
	Coordinates
	Status `json:"status"`
}
